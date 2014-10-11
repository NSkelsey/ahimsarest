package ahimsarest

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/NSkelsey/ahimsarest/ahimsadb"
	"github.com/gorilla/mux"
)

var db *sql.DB

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

	bltn, err := ahimsadb.GetJsonBltn(db, txid)
	if err == sql.ErrNoRows {
		http.Error(w, "Bulletin does not exist", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJson(w, bltn)

}

// returns the http handler initialized with the api's routes
func Handler() http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/bulletin/{txid:([a-f]|[0-9]){64}}", BulletinHandler)

	return r
}
