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

func setupDb() *sql.DB {
	dbpath := os.Getenv("GOPATH") + "/src/github.com/NSkelsey/ahimsarest/test.db"
	dbpath = filepath.Clean(dbpath)
	var err error
	db, err = ahimsadb.LoadDb(dbpath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func newTestServer() *httptest.Server {
	handler := ahimsarest.Handler()
	ts := httptest.NewServer(handler)
	return ts
}

var statusCodeTests = []struct {
	endpoint   string
	args       []interface{}
	statuscode int
}{
	// Test /bulletin/
	{"/bulletin/34324", nil, 404},
	{"/bulletin/deadbeef2ffc35dcc191acca037bed1defb0cf4df19555320502766c05041a62",
		nil,
		404,
	},
	{"/bulletin/f7800712c20377c2d29680c1aecf2331d6f80f5a44510d30ceb2e30fd5dafdcf",
		nil,
		200,
	},
	{"/bulletin/b0a1ba6e40d8f35aac526eecbc05d82b2a6d3c8d6a316627f593cbe592a777be",
		nil,
		451,
	},
	// Test /blacklist
	{"/blacklist", nil, 200},
}

// Runs a series of tests to assert the api is returning the correct status codes.
func TestStatusCodes(t *testing.T) {

	setupDb()
	ts := newTestServer()
	defer ts.Close()

	for _, testCase := range statusCodeTests {
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

var responseTests = []struct {
	endpoint string
	body     string
}{
	{
		endpoint: "/blacklist",
		body:     "[{\"txid\":\"b0a1ba6e40d8f35aac526eecbc05d82b2a6d3c8d6a316627f593cbe592a777be\",\"reason\":\"The Beatles are slanderous.\"}]",
	},
}

// Executes tests to verify that the json returned at an endpoint is correct
func TestResponses(t *testing.T) {
	setupDb()
	ts := newTestServer()
	defer ts.Close()

	for _, testCase := range responseTests {
		t.Logf("Endpoint: %s", testCase.endpoint)
		url := ts.URL + testCase.endpoint
		res, err := http.Get(url)
		if err != nil {
			t.Logf("%v", res)
			t.Error(err)
		}
		if res.StatusCode != 200 {
			t.Errorf("Responded with a %d!", res.StatusCode)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Error(err)
		}
		if string(body) != testCase.body {
			t.Logf("Responded with body:==========\n%s\nWanted:=========\n%s\n", body, testCase.body)
			t.Errorf("Bad json response")
		}
	}
}
