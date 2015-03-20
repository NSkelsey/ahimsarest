package main

import (
	"log"
	"net/http"

	"github.com/NSkelsey/ahimsadb"
	"github.com/NSkelsey/ahimsarest"
)

func main() {

	dbpath := "/home/ubuntu/.ahimsa/pubrecord.db"
	curdir := "/home/ubuntu/ahimsa-ang/"

	db, err := ahimsadb.LoadDB(dbpath)
	if err != nil {
		log.Fatal(err)
	}

	prefix := "/api/"
	api := ahimsarest.Handler(prefix, db)

	mux := http.NewServeMux()
	mux.Handle(prefix, api)
	mux.Handle("/", http.FileServer(http.Dir(curdir)))
	host := "0.0.0.0:1055"

	log.Printf("webserver listening at %s.\n", host)
	log.Fatal(http.ListenAndServe(host, mux))
}
