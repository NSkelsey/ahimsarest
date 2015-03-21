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
	staticpath = flag.String("statpath", "./", "The path to the static files to serve")
)

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
	log.Fatal(http.ListenAndServe(*host, mux))
}
