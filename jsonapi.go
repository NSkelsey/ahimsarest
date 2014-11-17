package ahimsarest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NSkelsey/ahimsadb"
	"github.com/NSkelsey/ahimsarest/ahimsajson"
	"github.com/NSkelsey/protocol/ahimsa"
	"github.com/gorilla/mux"
)

var (
	processStart time.Time = time.Now()
)

func writeJson(w http.ResponseWriter, m interface{}) {

	bytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, "Failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func BulletinHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		txid, _ := mux.Vars(request)["txid"]
		bltn, err := db.GetJsonBltn(txid)
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
}

// Handles requests for individual Blocks
func BlockHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		hash, _ := mux.Vars(request)["hash"]

		blockH, err := db.GetJsonBlock(hash)
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
}

// Handles a request for information about an individual author
func AuthorHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		addr, _ := mux.Vars(request)["addr"]

		authorJson, err := db.GetJsonAuthor(addr)
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
}

// Handles serving the blacklist contents over http. If the black list is empty
// it serves an empty list.
func BlacklistHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		blacklist, err := db.GetJsonBlacklist()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJson(w, blacklist)
	}
}

// Handles serving a bulletin board.
func BoardHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		boardstr, _ := mux.Vars(request)["board"]

		board, err := db.GetWholeBoard(boardstr)
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), 404)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, board)
	}
}

// Returns all bulletins under the board that has no name! Since board is an
// optional field you don't actually have to specify one. If that's the case
// then your bulletins will just have a NULL value in the board column
func NoBoardHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		board, err := db.GetWholeBoard("")
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), 404)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, board)
	}
}

// Returns the summaries of every board in the system sorted in lexicographic order.
func AllBoardsHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		boards, err := db.GetAllBoards()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, boards)
	}
}

// Returns all of the authors in the public record sorted in alphabetical order
func AllAuthorsHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		authors, err := db.GetAllAuthors()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, authors)
	}
}

// Returns all of the bulletins seen within the last 6 blocks.
func RecentHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		bltns, err := db.GetRecentConf(6)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, bltns)
	}
}

// Returns all of the unconfirmed bulletins ordered by reported time.
func UnconfirmedHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		bltns, err := db.GetUnconfirmed()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, bltns)
	}
}

// Returns all of the block summaries for a given day.
func BlockDayHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		datestr := mux.Vars(request)["day"]

		// convert into UTC then do lookups within range
		layout := "02-01-2006"
		date, err := time.Parse(layout, datestr)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		blocks, err := db.GetBlocksByDay(date)
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), 404)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJson(w, blocks)
	}
}

// Handles the round trip to ahimsadb to get DB status. In the future
// this could look up the status of other processes that are running
// on the machine and report their status as well.
func StatusHandler(db *ahimsadb.PublicRecord) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		latestBlk, latestBltn, err := db.LatestBlkAndBltn()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		status := &ahimsajson.Status{
			Version:    ahimsa.Version,
			AppStart:   processStart.Unix(),
			LatestBlk:  latestBlk,
			LatestBltn: latestBltn,
		}

		writeJson(w, status)
	}
}

// returns the http handler initialized with the api's routes
func Handler(db *ahimsadb.PublicRecord) http.Handler {

	r := mux.NewRouter()
	sha2re := "([a-f]|[A-F]|[0-9]){64}"
	addrgex := "([a-z]|[A-Z]|[0-9]){30,35}"
	// Since the board's path could be percent encoded we give it 3x wiggle room
	// since a single byte in percent encoding is %EE.
	boardre := ".{1,90}"

	// A single day follows this format: DD-MM-YY
	dayre := `[0-9]{1,2}-[0-9]{1,2}-[0-9]{4}`

	// Item handlers
	r.HandleFunc(fmt.Sprintf("/bulletin/{txid:%s}", sha2re), BulletinHandler(db))
	r.HandleFunc(fmt.Sprintf("/author/{addr:%s}", addrgex), AuthorHandler(db))
	r.HandleFunc(fmt.Sprintf("/block/{hash:%s}", sha2re), BlockHandler(db))
	r.HandleFunc(fmt.Sprintf("/board/{board:%s}", boardre), BoardHandler(db))
	r.HandleFunc("/blacklist", BlacklistHandler(db))
	r.HandleFunc("/noboard", NoBoardHandler(db))

	// Aggregate handlers
	r.HandleFunc("/boards", AllBoardsHandler(db))
	r.HandleFunc("/recent", RecentHandler(db))
	r.HandleFunc("/unconfirmed", UnconfirmedHandler(db))
	r.HandleFunc("/authors", AllAuthorsHandler(db))
	r.HandleFunc(fmt.Sprintf("/blocks/{day:%s}", dayre), BlockDayHandler(db))

	// Meta handlers
	r.HandleFunc("/status", StatusHandler(db))

	return r
}
