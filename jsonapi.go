package ahimsarest

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	if err == ahimsadb.ErrBltnCensored {
		http.Error(w, err.Error(), 451)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJson(w, bltn)

}

// Handles requests for individual Blocks
func BlockHandler(w http.ResponseWriter, request *http.Request) {

	hash, _ := mux.Vars(request)["hash"]

	// TODO implement
	blockH, err := ahimsadb.GetJsonBlock(db, hash)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJson(w, blockH)
}

// Handles a request for information about an individual author
func AuthorHandler(w http.ResponseWriter, request *http.Request) {

	addr, _ := mux.Vars(request)["addr"]

	// TODO write testcases
	authorJson, err := ahimsadb.GetJsonAuthor(db, addr)
	if err == sql.ErrNoRows {
		http.Error(w, "Author does not exist", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJson(w, authorJson)
}

// Handles serving the blacklist contents over http. If the black list is empty
// it serves an empty list.
func BlacklistHandler(w http.ResponseWriter, request *http.Request) {
	blacklist, err := ahimsadb.GetJsonBlacklist(db)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJson(w, blacklist)
}

// returns the http handler initialized with the api's routes
func Handler() http.Handler {

	r := mux.NewRouter()
	sha2re := "([a-f]|(A-F)|[0-9]){64}"
	addrgex := "([a-f]|(A-F)|[0-9]){30,35}"

	// Item handlers
	r.HandleFunc(fmt.Sprintf("/bulletin/{txid:%s}", sha2re), BulletinHandler)
	r.HandleFunc(fmt.Sprintf("/author/{addr:%s}", addrgex), AuthorHandler)

	// Aggregate handlers
	r.HandleFunc(fmt.Sprintf("/block/{hash:%s}", sha2re), BlockHandler)

	r.HandleFunc("/blacklist", BlacklistHandler)

	return r
}
