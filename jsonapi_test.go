package ahimsarest_test

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/NSkelsey/ahimsarest"
	"github.com/NSkelsey/ahimsarest/ahimsadb"
)

var db *sql.DB

var cmdtests = []struct {
	endpoint   string
	args       []interface{}
	statuscode int
}{
	{"/bulletin/34324", nil, 404},
	{"/bulletin/deadbeef2ffc35dcc191acca037bed1defb0cf4df19555320502766c05041a62",
		nil,
		404,
	},
	{"/bulletin/3b3415ca2ffc35dcc191acca037bed1defb0cf4df19555320502766c05041a62",
		nil,
		200,
	},
}

// Runs a series of end to end tests that assert the functionality of the api's
// endpoints.
func TestJsonApi(t *testing.T) {

	// TODO use test db setup
	dbpath := os.Getenv("GOPATH") + "/src/github.com/NSkelsey/ahimsarest/test.db"
	dbpath = filepath.Clean(dbpath)
	var err error
	db, err = ahimsadb.LoadDb(dbpath)
	if err != nil {
		log.Fatal(err)
	}

	handler := ahimsarest.Handler()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	for _, testCase := range cmdtests {
		url := ts.URL + testCase.endpoint
		res, err := http.Get(url)
		if err != nil {
			t.Logf("%v", res)
			t.Error(err)
		}

		if res.StatusCode != testCase.statuscode {
			t.Logf("Expected: %d, Recieved: %d\n", testCase.statuscode, res.StatusCode)
			body, _ := ioutil.ReadAll(res.Body)
			res.Body.Close()
			t.Logf("Resp Body:\n%s", body)
			t.Errorf("Endpoint: %s status code does not meet expectations\n", testCase.endpoint)
		}
	}
}
