package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/NSkelsey/ahimsarest"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/pubrecdb"
)

var (
	host       = flag.String("host", "localhost:1055", "The ip and port for the server to listen on")
	staticpath = flag.String("staticpath", "./", "The path to the static files to serve")
	verbose    = flag.Bool("verbose", false, "Logs the output of every request")
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()

	nodedir := filepath.Join(btcutil.AppDataDir("ombudscore", false), "node")
	dbpath := filepath.Join(nodedir, "pubrecord.db")

	db, err := pubrecdb.LoadDB(dbpath)
	if err != nil {
		log.Fatal(err)
	}

	prefix := "/api/"
	api := ahimsarest.Handler(prefix, db)

	mux := http.NewServeMux()

	mux.Handle(prefix, api)
	mux.Handle("/", http.FileServer(http.Dir(*staticpath)))

	log.Printf("webserver listening at %s.\n", *host)

	if *verbose {
		logger := Log(mux)
		log.Fatal(http.ListenAndServe(*host, logger))
	} else {
		log.Fatal(http.ListenAndServe(*host, mux))
	}
}
