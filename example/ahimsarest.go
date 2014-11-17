package main

import (
	"log"
	"net/http"

	"github.com/NSkelsey/ahimsadb"
	"github.com/NSkelsey/ahimsarest"
)

func main() {

	dbpath := "/home/ubuntu/.ahimsa/pubrecord.db"

	var err error
	db, err := ahimsadb.LoadDb(dbpath)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", ahimsarest.Handler(db))
	host := "0.0.0.0:1054"
	log.Printf("ahimsarest listening at %s.\n", host)
	http.ListenAndServe(host, nil)

}
