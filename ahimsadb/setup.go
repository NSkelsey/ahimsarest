package ahimsadb

import (
	"database/sql"
	"os"
	"path/filepath"
)

func SetupTestDB() (*sql.DB, error) {

	dbpath := os.Getenv("GOPATH") + "/src/github.com/NSkelsey/ahimsarest/test.db"
	dbpath = filepath.Clean(dbpath)
	var err error
	db, err := LoadDb(dbpath)
	if err != nil {
		return nil, err
	}
	return db, nil
}
