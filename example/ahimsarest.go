package main

import (
	"log"
	"net/http"

	"github.com/NSkelsey/ahimsarest"
	"github.com/NSkelsey/ahimsarest/ahimsadb"
)

func main() {

	dbpath := "/home/ubuntu/gocode/src/github.com/NSkelsey/ahimsarest/test.db"

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
