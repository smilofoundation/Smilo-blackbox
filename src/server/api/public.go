package api

import (
	"net/http"
	"github.com/gorilla/mux"
)

func GetPartyInfo(w http.ResponseWriter, r *http.Request) {

}

func Push(w http.ResponseWriter, r *http.Request) {

}

func Resend(w http.ResponseWriter, r *http.Request) {

}

func Delete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
    epk := r.PostForm.Get("Encoded public key")
    w.Write([]byte(epk))
}

func TransactionGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Write([]byte(params["hash"]))
}

func TransactionDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Write([]byte(params["key"]))
}