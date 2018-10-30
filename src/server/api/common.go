package api

import "net/http"

func GetVersion(w http.ResponseWriter, r *http.Request) {

}
func Upcheck(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("I'm up!"))
}