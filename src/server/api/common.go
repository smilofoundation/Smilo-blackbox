package api

import "net/http"
const BlackBoxVersion = "Smilo Black Box 0.1.0"

func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte([]byte(BlackBoxVersion)))
}
func Upcheck(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("I'm up!"))
}