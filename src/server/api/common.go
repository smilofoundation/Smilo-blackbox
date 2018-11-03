package api

import "net/http"
const BlackBoxVersion = "Smilo Black Box 0.1.0"
const UpcheckMessage = "I'm up!"

// Request path "/version", response plain text version ID
func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(BlackBoxVersion))
}

// Request path "/upcheck", response plain text upcheck message.
func Upcheck(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(UpcheckMessage))
}

// Request path "/api", response json rest api spec.
func Api(w http.ResponseWriter, r *http.Request) {

}