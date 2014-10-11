package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/NSkelsey/ahimsarest"
	"github.com/NSkelsey/ahimsarest/ahimsadb"
)

var db *sql.DB

func main() {

	dbpath := "/home/ubuntu/.ahimsa/pubrecord.db"

	var err error
	db, err = ahimsadb.LoadDb(dbpath)
	if err != nil {
		log.Fatal(err)
	}
	// write items funcs first

	http.Handle("/", ahimsarest.Handler())
	log.Println("listening")
	http.ListenAndServe(":1054", nil)

}
