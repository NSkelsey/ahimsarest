package main

import (
	"encoding/json"
	"net/http"

	"github.com/NSkelsey/ahimsarest/ahimsajson"
	"github.com/gorilla/mux"
)

func writeJson(w http.ResponseWriter, m interface{}) {

	bytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, "Failed", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func BulletinHandler(w http.ResponseWriter, request *http.Request) {

	txid, _ := mux.Vars(request)["txid"]

	// add DB pull to create json object.

	bltn := ahimsajson.Bulletin{
		Txid:    txid,
		Message: "halp halp I am being repressed",
		Board:   "ahimsa-dev",
	}

	wrapWrite(bltn, w)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/bulletin/{txid:([a-f]|[0-9]){64}}", BulletinHandler)
	// write items funcs first

	http.Handle("/", r)
	http.ListenAndServe(":8083", nil)

}
