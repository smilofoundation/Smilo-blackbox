package server

import (
	"os"
	"net/http"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
	"Smilo-blackbox/src/server/api"
	"net"
	"github.com/gorilla/mux"
)
var privateAPI             *mux.Router
var publicAPI              *mux.Router

var	sockPath = os.TempDir()+"blackbox.sock"

type status struct {
	httpServer bool
	unixServer bool
}

var serverStatus = status{false, false}

func NewServer(Port string) (*http.Server,*http.Server) {
	publicAPI, privateAPI = InitRouting()

	return &http.Server{
		Addr:    ":" + Port,
		Handler: publicAPI,
	},
		&http.Server{
			Handler: privateAPI,
		}

}


func StartServer(Port string, sockFilePath string) {

	if sockFilePath != "" {
		sockPath = sockFilePath
	}
	os.Remove(sockPath)

	glog.Info("Starting server")
	pub, priv := NewServer(Port)
	glog.Info("Server starting --> " + Port)
	sock, _ := net.Listen("unix", sockPath)
	glog.Info("Unix Domain Socket Up --> " + sockPath)
    go func() {
		err := gracehttp.Serve(
			pub)
		if err != nil {
			glog.Error("Error: %v", err)
			os.Exit(1)
		}
	}()
	err := priv.Serve(sock)
	if err != nil {
		glog.Error("Error: %v", err)
		os.Exit(1)
	}
}

func InitRouting() (*mux.Router, *mux.Router) {

	publicAPI := mux.NewRouter()
	privateAPI := mux.NewRouter()

	publicAPI.HandleFunc("/version", api.GetVersion).Methods("GET")
	publicAPI.HandleFunc("/upcheck", api.Upcheck).Methods("GET")

	publicAPI.HandleFunc("/push", api.Push).Methods("POST")
	publicAPI.HandleFunc("/resend", api.Resend).Methods("POST")
	publicAPI.HandleFunc("/partyinfo", api.GetPartyInfo).Methods("GET")
	publicAPI.HandleFunc("/delete", api.Delete).Methods("POST")
	publicAPI.HandleFunc("/transaction/{hash:.*}",api.TransactionGet).Methods("GET")
	publicAPI.HandleFunc("/transaction/{key:.*}",api.TransactionDelete).Methods("DELETE")

	privateAPI.HandleFunc("/version", api.GetVersion).Methods("GET")
	privateAPI.HandleFunc("/upcheck", api.Upcheck).Methods("GET")
	privateAPI.HandleFunc("/sendraw", api.SendRaw).Methods("POST")
	privateAPI.HandleFunc("/receiveraw", api.ReceiveRaw).Methods("GET")

	return publicAPI, privateAPI
}
