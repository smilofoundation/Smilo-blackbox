// Copyright 2019 The Smilo-blackbox Authors
// This file is part of the Smilo-blackbox library.
//
// The Smilo-blackbox library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Smilo-blackbox library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Smilo-blackbox library. If not, see <http://www.gnu.org/licenses/>.

package server

import (
	"net"
	"net/http"
	"os"
	"strconv"

	"Smilo-blackbox/src/server/api"

	"crypto/tls"
	"path"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"

	"sync"
	"time"

	"github.com/asdine/storm"
	"github.com/tidwall/buntdb"

	"Smilo-blackbox/src/data"
	"Smilo-blackbox/src/server/config"
	"Smilo-blackbox/src/utils"
)

var (
	msgC       = make(chan Message)
	msgCmutex  = &sync.Mutex{}
	privateAPI *mux.Router
	publicAPI  *mux.Router

	//StormDBPeers is the main object for peers db
	StormDBPeers *storm.DB

	//serverStatus               = status{false, false}
	log *logrus.Entry = logrus.WithFields(logrus.Fields{
		"app":     "blackbox",
		"package": "server",
	})

	//DefaultExpirationTime is the default expiration time used on the database
	DefaultExpirationTime = &buntdb.SetOptions{Expires: false} // never expire

	//serverURL             string

	//PUBLIC_SERVER_READ_TIMEOUT_STR will be used to hold env var
	PUBLIC_SERVER_READ_TIMEOUT_STR = os.Getenv("PUBLIC_SERVER_READ_TIMEOUT")

	//PUBLIC_SERVER_WRITE_TIMEOUT_STR will be used to hold env var
	PUBLIC_SERVER_WRITE_TIMEOUT_STR = os.Getenv("PUBLIC_SERVER_WRITE_TIMEOUT")

	//PRIVATE_SERVER_READ_TIMEOUT_STR will be used to hold env var
	PRIVATE_SERVER_READ_TIMEOUT_STR = os.Getenv("PRIVATE_SERVER_READ_TIMEOUT")

	//PRIVATE_SERVER_WRITE_TIMEOUT_STR will be used to hold env var
	PRIVATE_SERVER_WRITE_TIMEOUT_STR = os.Getenv("PRIVATE_SERVER_WRITE_TIMEOUT")
	//PUBLIC_SERVER_READ_TIMEOUT will be used to hold env var
	PUBLIC_SERVER_READ_TIMEOUT = 120
	//PUBLIC_SERVER_WRITE_TIMEOUT will be used to hold env var
	PUBLIC_SERVER_WRITE_TIMEOUT = 120
	//PRIVATE_SERVER_READ_TIMEOUT   will be used to hold env var
	PRIVATE_SERVER_READ_TIMEOUT = 60
	//PRIVATE_SERVER_WRITE_TIMEOUT  will be used to hold env var
	PRIVATE_SERVER_WRITE_TIMEOUT = 60
)

func init() {

}

func setOSEnvInt(v, field string, defaultVal int) (r int) {
	var err error
	r, err = strconv.Atoi(v)
	if err != nil {
		log.WithError(err).Warnf("Going to use default field %s, defaultVal: %d", field, defaultVal)
	} else {
		return defaultVal
	}
	return r
}

func initServer() {
	PUBLIC_SERVER_READ_TIMEOUT = setOSEnvInt(PUBLIC_SERVER_READ_TIMEOUT_STR, "PUBLIC_SERVER_READ_TIMEOUT", PUBLIC_SERVER_READ_TIMEOUT)
	PUBLIC_SERVER_WRITE_TIMEOUT = setOSEnvInt(PUBLIC_SERVER_WRITE_TIMEOUT_STR, "PUBLIC_SERVER_WRITE_TIMEOUT", PUBLIC_SERVER_WRITE_TIMEOUT)
	PRIVATE_SERVER_READ_TIMEOUT = setOSEnvInt(PRIVATE_SERVER_READ_TIMEOUT_STR, "PRIVATE_SERVER_READ_TIMEOUT", PRIVATE_SERVER_READ_TIMEOUT)
	PRIVATE_SERVER_WRITE_TIMEOUT = setOSEnvInt(PRIVATE_SERVER_WRITE_TIMEOUT_STR, "PRIVATE_SERVER_WRITE_TIMEOUT", PRIVATE_SERVER_WRITE_TIMEOUT)

	finalPath := utils.BuildFilename(config.PeersDBFile.Value)
	_, err := os.Create(finalPath)
	if err != nil {
		log.Fatalf("Failed to start StormDBPeers file at %s", config.Socket.Value)
	}

	StormDBPeers, err = storm.Open(finalPath)
	if err != nil {
		defer func() {
			err := StormDBPeers.Close()
			if err != nil {
				log.WithError(err).Error("Could not StormDBPeers.Close")
			}
		}()
		log.WithError(err).Error("Could not open StormDBPeers")
		os.Exit(1)
	}

}

// SetLogger set the logger
func SetLogger(loggers *logrus.Entry) {
	log = loggers.WithFields(log.Data)

	filenameHook := filename.NewHook()

	logrus.AddHook(filenameHook)

}

//NewServer will create a new http server instance -- pub and private
func NewServer(Port string) (*http.Server, *http.Server) {
	publicAPI, privateAPI = InitRouting()

	return &http.Server{
			Addr:         ":" + Port,
			Handler:      publicAPI,
			ReadTimeout:  time.Duration(PUBLIC_SERVER_READ_TIMEOUT) * time.Second,
			WriteTimeout: time.Duration(PUBLIC_SERVER_WRITE_TIMEOUT) * time.Second,
		},
		&http.Server{
			Handler:      privateAPI,
			ReadTimeout:  time.Duration(PRIVATE_SERVER_READ_TIMEOUT) * time.Second,
			WriteTimeout: time.Duration(PRIVATE_SERVER_WRITE_TIMEOUT) * time.Second,
		}

}

//StartServer will start the server
func StartServer() {
	port, isTLS, workDir := config.Port.Value, config.IsTLS.Destination, config.WorkDir.Value
	initServer()
	log.Info("Starting server")
	pub, priv := NewServer(port)

	log.Info("Server starting --> " + port)
	data.Start()

	if *isTLS {
		log.Info("Will start TLS Mode")
		servCert := config.ServCert.Value
		servKey := config.ServKey.Value

		if (len(servCert) != len(servKey)) || (len(servCert) <= 0) {
			log.Fatalf("Please provide server certificate and key for TLS %s %s %d ", servKey, servCert, len(servCert))
		}

		certFile := path.Join(workDir, servCert)
		keyFile := path.Join(workDir, servKey)

		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			log.Error(err)
			os.Exit(1)
		} else if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			log.Error(err)
			os.Exit(1)
		}

		go func() {
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				log.Fatalf("Error loading cert: %v", err)
				os.Exit(1)
			}
			pub.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
			err = gracehttp.Serve(pub)
			if err != nil {
				log.Fatalf("Error starting server with TLS: %v", err)
				os.Exit(1)
			}
		}()

	} else {
		go func() {
			err := gracehttp.Serve(pub)
			if err != nil {
				log.Fatalf("Error starting API server: %v", err)
				os.Exit(0)
			}
		}()
	}

	finalPath := utils.BuildFilename(config.Socket.Value)
	err := os.Remove(finalPath)
	if err != nil {
		log.WithError(err).Error("Could not os.Remove(finalPath), 1")
	}

	time.Sleep(1 * time.Second)
	err = os.MkdirAll(finalPath, os.FileMode(0755))
	if err != nil {
		log.Fatalf("Failed to start IPC Server at %s", config.Socket.Value)
	}

	err = os.Remove(finalPath)
	if err != nil {
		log.WithError(err).Error("Could not os.Remove(finalPath), 2")
	}

	defer func() {
		err = os.Remove(finalPath)
		if err != nil {
			log.WithError(err).Error("Could not os.Remove(finalPath), 3")
		}
	}()

	log.Info("Starting IPC Server at, ", config.Socket.Value)
	go func() {
		sock, err := net.Listen("unix", finalPath)
		if err != nil {
			log.Fatalf("Failed to start IPC Server at %s", config.Socket.Value)
		}
		err = os.Chmod(finalPath, 0600)
		if err != nil {
			log.WithError(err).Error("Could not os.Chmod(finalPath, 0600)")
		}

		err = priv.Serve(sock)
		if err != nil {
			log.Errorf("Error: %v", err)
			os.Exit(1)
		}
	}()
}

//InitRouting will init routing
func InitRouting() (*mux.Router, *mux.Router) {

	publicAPI := mux.NewRouter()
	privateAPI := mux.NewRouter()

	publicAPI.HandleFunc("/upcheck", api.Upcheck).Methods("GET")
	publicAPI.HandleFunc("/version", api.GetVersion).Methods("GET")
	publicAPI.HandleFunc("/push", api.Push).Methods("POST")
	publicAPI.HandleFunc("/resend", api.Resend).Methods("POST")
	publicAPI.HandleFunc("/storeraw", api.StoreRaw).Methods("POST")
	publicAPI.HandleFunc("/partyinfo", api.GetPartyInfo).Methods("POST")
	publicAPI.NotFoundHandler = http.HandlerFunc(api.UnknownRequest)

	// Restrict to IPC
	privateAPI.HandleFunc("/upcheck", api.Upcheck).Methods("GET")
	privateAPI.HandleFunc("/version", api.GetVersion).Methods("GET")
	privateAPI.HandleFunc("/send", api.Send).Methods("POST")
	privateAPI.HandleFunc("/sendraw", api.SendRaw).Methods("POST")
	privateAPI.HandleFunc("/receive", api.Receive).Methods("GET")
	privateAPI.HandleFunc("/sendsignedtx", api.SendSignedTx).Methods("POST")
	privateAPI.HandleFunc("/receiveraw", api.ReceiveRaw).Methods("GET")
	privateAPI.HandleFunc("/delete", api.Delete).Methods("POST")

	privateAPI.HandleFunc("/transaction/{hash:.*}", api.TransactionGet).Methods("GET")

	publicAPI.HandleFunc("/transaction/{key:.*}", api.TransactionDelete).Methods("DELETE")
	publicAPI.HandleFunc("/config/peers", api.ConfigPeersPut).Methods("PUT")
	publicAPI.HandleFunc("/config/peers/{index:.*}", api.ConfigPeersGet).Methods("GET")
	publicAPI.HandleFunc("/metrics", api.Metrics).Methods("GET")

	return publicAPI, privateAPI
}
