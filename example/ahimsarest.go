package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/NSkelsey/ahimsarest"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/pubrecdb"
)

func main() {

	dbpath := filepath.Join(btcutil.AppDataDir("ombfullnode", false), "pubrecord.db")

	var err error
	db, err := pubrecdb.LoadDB(dbpath)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", ahimsarest.Handler("", db))
	host := "0.0.0.0:1054"
	log.Printf("web-api listening at %s.\n", host)
	http.ListenAndServe(host, nil)

}
