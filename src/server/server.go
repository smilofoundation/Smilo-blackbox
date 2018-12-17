package server

import (
	"net"
	"net/http"
	"os"

	"Smilo-blackbox/src/server/api"

	"crypto/tls"
	"path"
	"path/filepath"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"

	"Smilo-blackbox/src/server/config"
)

var (
	privateAPI *mux.Router
	publicAPI  *mux.Router

	serverStatus               = status{false, false}
	log          *logrus.Entry = logrus.WithField("package", "server")
)

type status struct {
	httpServer bool
	unixServer bool
}

// SetLogger set the logger
func SetLogger(loggers *logrus.Entry) {
	log = loggers.WithFields(log.Data)

	filenameHook := filename.NewHook()

	logrus.AddHook(filenameHook)

}

func NewServer(Port string) (*http.Server, *http.Server) {
	publicAPI, privateAPI = InitRouting()

	return &http.Server{
			Addr:    ":" + Port,
			Handler: publicAPI,
		},
		&http.Server{
			Handler: privateAPI,
		}

}

func StartServer() {
	port, isTLS, workDir := config.Port.Value, config.IsTLS.Value, config.WorkDir.Value

	log.Info("Starting server")
	pub, priv := NewServer(port)
	log.Info("Server starting --> " + port)

	if isTLS != "" {
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
			err := gracehttp.Serve(
				pub)
			if err != nil {
				log.Fatalf("Error starting server: %v", err)
				os.Exit(1)
			}
		}()
	}

	socketFile := config.Socket.Value
	err := os.MkdirAll(filepath.Join(workDir, socketFile), os.FileMode(0755))
	if err != nil {
		log.Fatalf("Failed to start IPC Server at %s", socketFile)
	}

	os.Remove(socketFile)

	defer func() {
		os.Remove(socketFile)
	}()

	sock, err := net.Listen("unix", socketFile)
	if err != nil {
		log.Fatalf("Failed to start IPC Server at %s", socketFile)
	}
	os.Chmod(socketFile, 0600)

	err = priv.Serve(sock)
	if err != nil {
		log.Error("Error: %v", err)
		os.Exit(1)
	}
}

func InitRouting() (*mux.Router, *mux.Router) {

	publicAPI := mux.NewRouter()
	privateAPI := mux.NewRouter()

	publicAPI.HandleFunc("/upcheck", api.Upcheck).Methods("GET")
	publicAPI.HandleFunc("/version", api.GetVersion).Methods("GET")
	publicAPI.HandleFunc("/push", api.Push).Methods("POST")
	publicAPI.HandleFunc("/resend", api.Resend).Methods("POST")
	publicAPI.HandleFunc("/partyinfo", api.GetPartyInfo).Methods("GET")

	// Restrict to IPC
	privateAPI.HandleFunc("/upcheck", api.Upcheck).Methods("GET")
	privateAPI.HandleFunc("/version", api.GetVersion).Methods("GET")
	privateAPI.HandleFunc("/send", api.Send).Methods("POST")
	privateAPI.HandleFunc("/sendraw", api.SendRaw).Methods("POST")
	privateAPI.HandleFunc("/receive", api.Receive).Methods("GET")
	privateAPI.HandleFunc("/receiveraw", api.ReceiveRaw).Methods("GET")
	privateAPI.HandleFunc("/delete", api.Delete).Methods("POST")

	privateAPI.HandleFunc("/transaction/{hash:.*}", api.TransactionGet).Methods("GET")

	publicAPI.HandleFunc("/transaction/{key:.*}", api.TransactionDelete).Methods("DELETE")
	publicAPI.HandleFunc("/config/peers", api.ConfigPeersPut).Methods("PUT")
	publicAPI.HandleFunc("/config/peers/{index:.*}", api.ConfigPeersGet).Methods("GET")
	publicAPI.HandleFunc("/metrics", api.Metrics).Methods("GET")

	return publicAPI, privateAPI
}
