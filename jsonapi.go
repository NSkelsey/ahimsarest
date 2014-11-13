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

	blockH, err := ahimsadb.GetJsonBlock(db, hash)
	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJson(w, blockH)
}

// Handles a request for information about an individual author
func AuthorHandler(w http.ResponseWriter, request *http.Request) {

	addr, _ := mux.Vars(request)["addr"]

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

// Handles serving a bulletin board.
func BoardHandler(w http.ResponseWriter, request *http.Request) {
	boardstr, _ := mux.Vars(request)["board"]

	board, err := ahimsadb.GetWholeBoard(db, boardstr)
	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), 404)
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJson(w, board)
}

// Returns all bulletins under the board that has no name! Since board is an
// optional field you don't actually have to specify one. If that's the case
// then your bulletins will just have a NULL value in the board column
func NoBoardHandler(w http.ResponseWriter, request *http.Request) {

	board, err := ahimsadb.GetWholeBoard(db, "")
	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), 404)
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	writeJson(w, board)
}

// Returns the summaries of every board in the system sorted in lexicographic order.
func AllBoardsHandler(w http.ResponseWriter, request *http.Request) {

	boards, err := ahimsadb.GetAllBoards(db)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	writeJson(w, boards)
}

// Returns all of the bulletins seen within the last 6 blocks.
func RecentHandler(w http.ResponseWriter, request *http.Request) {

	bltns, err := ahimsadb.GetRecentConf(db, 6)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	writeJson(w, bltns)
}

// Returns all of the unconfirmed bulletins ordered by reported time.
func UnconfirmedHandler(w http.ResponseWriter, request *http.Request) {

	bltns, err := ahimsadb.GetUnconfirmed(db)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	writeJson(w, bltns)
}

// Returns all of the block summaries for a given day.
func BlockDayHandler(w http.ResponseWriter, request *http.Request) {

	// TODO write testcases
	day := int(mux.Vars.request["day"])
	month := int(mux.Vars.request["month"])
	year := int(mux.Vars.request["year"])

	logger.Println(day, month, year)

	// TODO test + convert into UTC then do lookups withinin range

	//writeJson(w, blocks)
}

// returns the http handler initialized with the api's routes
func Handler() http.Handler {

	r := mux.NewRouter()
	sha2re := "([a-f]|[A-F]|[0-9]){64}"
	addrgex := "([a-z]|[A-Z]|[0-9]){30,35}"
	// Since the board's path could be percent encoded we give it 3x wiggle room
	// since a single byte in percent encoding is %EE.
	boardre := ".{1,90}"

	// A single day follows this format: DD-MM-YY
	dayre := "{day:[0-9]{2}}-{month:[0-9]{2}}-{year:[0-9]{4}}"

	// Item handlers
	r.HandleFunc(fmt.Sprintf("/bulletin/{txid:%s}", sha2re), BulletinHandler)
	r.HandleFunc(fmt.Sprintf("/author/{addr:%s}", addrgex), AuthorHandler)
	r.HandleFunc(fmt.Sprintf("/block/{hash:%s}", sha2re), BlockHandler)
	r.HandleFunc(fmt.Sprintf("/board/{board:%s}", boardre), BoardHandler)
	r.HandleFunc("/blacklist", BlacklistHandler)
	r.HandleFunc("/noboard", NoBoardHandler)

	// Aggregate handlers
	r.HandleFunc("/boards", AllBoardsHandler)
	r.HandleFunc("/recent", RecentHandler)
	r.HandleFunc("/unconfirmed", UnconfirmedHandler)
	r.HandleFunc(fmt.Sprintf("/blocks/%s", dayre), BlockDayHandler)

	return r
}
