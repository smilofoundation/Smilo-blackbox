package server

import (
	"os"
	"net/http"
	"github.com/orisano/uds"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/pat"
	"github.com/golang/glog"
	"Smilo-blackbox/src/server/api"
	"sync"
)
var privateAPI             *pat.Router
var publicAPI              *pat.Router

const (
	sockPath = "./sample.sock"
)

func NewServer(Port string) *http.Server {
	publicAPI, privateAPI = InitRouting()
	return &http.Server{
		Addr:    ":" + Port,
		Handler: publicAPI,
	}
}

//StartServer start and listen @server
func StartServer(Port string) {

	//gracehttp.SetLogger(logger)

	glog.Info("Starting server")
	s := NewServer(Port)
	glog.Info("Server starting --> " + Port)

	var wg sync.WaitGroup
	//enable graceful shutdown
	wg.Add(2)
	go func() {
		err := gracehttp.Serve(
			s,
		)
		if err != nil {
			glog.Error("Error: %v", err)
			os.Exit(1)
		}
        wg.Done()
	}()

	go func() {
		os.Remove(sockPath)
		err := uds.ListenAndServe(sockPath, privateAPI)
		if err != nil {
			glog.Error("Error: %v", err)
			os.Exit(1)
		}
		wg.Done()
	}()
	wg.Wait()

}

func InitRouting() (*pat.Router, *pat.Router) {

	publicAPI := pat.New()
	privateAPI := pat.New()

	publicAPI.HandleFunc("/version", api.GetVersion).Methods("GET")
	publicAPI.HandleFunc("/upcheck", api.Upcheck).Methods("GET")

	publicAPI.HandleFunc("/push", api.Push).Methods("POST")
	publicAPI.HandleFunc("/resend", api.Resend).Methods("POST")
	publicAPI.HandleFunc("/partyinfo", api.GetPartyInfo).Methods("GET")
	publicAPI.HandleFunc("/delete", api.Delete).Methods("POST")

	privateAPI.HandleFunc("/version", api.GetVersion).Methods("GET")
	privateAPI.HandleFunc("/upcheck", api.Upcheck).Methods("GET")
	privateAPI.HandleFunc("/sendraw", api.SendRaw).Methods("POST")
	privateAPI.HandleFunc("/receiveraw", api.ReceiveRaw).Methods("GET")

	return publicAPI, privateAPI
}
