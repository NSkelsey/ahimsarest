package ahimsarest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/NSkelsey/ahimsarest/ahimsadb"
)

func newTestServer(t *testing.T) *httptest.Server {
	_, err := ahimsadb.SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	handler := Handler()
	ts := httptest.NewServer(handler)
	return ts
}

// Assert proper response codes based off provided url paths.
var statusCodeTests = []struct {
	endpoint   string
	statuscode int
}{
	// Test /bulletin/
	{"/bulletin/34324", 404},
	{"/bulletin/deadbeef2ffc35dcc191acca037bed1defb0cf4df19555320502766c05041a62",
		404,
	},
	{"/bulletin/f7800712c20377c2d29680c1aecf2331d6f80f5a44510d30ceb2e30fd5dafdcf",
		200,
	},
	{"/bulletin/b0a1ba6e40d8f35aac526eecbc05d82b2a6d3c8d6a316627f593cbe592a777be",
		451,
	},
	{"/block/0000000000000000000000000000000000000000000000000000000000000000",
		404,
	},
	{"/block/ThisShouldNotMatch",
		404,
	},
	{"/blacklist", 200},
	{"/author/0000000000000000000000000000000000", 404},
	{"/board/this-One-Isnt-Real", 404},
	// Ensure that a utf-8 url-encoded board is reachable
	{"/board/%23%21~%2AEnc%20ded-bo%C3%84&%5C/%D3%81", 200},
	{"/noboard", 200},
}

// Runs a series of tests to assert the api is returning the correct status codes.
func TestStatusCodes(t *testing.T) {

	ts := newTestServer(t)
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

// These endpoints should return valid responses with actual content
var responseTests = []struct {
	endpoint string
	body     string
}{
	{
		endpoint: "/bulletin/f7800712c20377c2d29680c1aecf2331d6f80f5a44510d30ceb2e30fd5dafdcf",
		body:     `{"txid":"f7800712c20377c2d29680c1aecf2331d6f80f5a44510d30ceb2e30fd5dafdcf","board":"ahimsa-dev","author":"mraY7GWs4G65ZYqWoPEqwPxb6aMuufTq2c","msg":"Here comes the sun","timestamp":1413079355}`,
	},
	{
		endpoint: "/block/00000000777213b4fd7c5d5a71b9b52608356c4194203b1b63d1bb0e6141d17d",
		body:     `{"head":{"hash":"00000000777213b4fd7c5d5a71b9b52608356c4194203b1b63d1bb0e6141d17d","prevHash":"00000000f1cb8c224acdb0e1becbfa4218f1e13b0d4dbbce64d0a3c15d8bf55f","timestamp":1414813562,"height":305724,"numBltns":1},"bltns":[{"txid":"b0a1ba6e40d8f35aac526eecbc05d82b2a6d3c8d6a316627f593cbe592a777be","board":"ahimsa-dev","author":"mnPZBNTrLoCoSkAgSfKeeCujU3129PG6vn","msg":"","timestamp":1413216499,"blk":"00000000777213b4fd7c5d5a71b9b52608356c4194203b1b63d1bb0e6141d17d","blkTimestamp":1414813562,"bannedReason":"The Beatles are slanderous."}]}`,
	},
	{
		endpoint: "/blacklist",
		body:     `[{"txid":"b0a1ba6e40d8f35aac526eecbc05d82b2a6d3c8d6a316627f593cbe592a777be","reason":"The Beatles are slanderous."}]`,
	},
	{
		endpoint: "/author/miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH",
		body:     `{"author":{"addr":"miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH","numBltns":2,"firstBlkTs":1414017952},"bltns":[{"txid":"2963cc35727f4e2c2bd4186e4550fe82b204e446ff7096b425f236264e05c7c6","board":"ahimsa-dev","author":"miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH","msg":"pier and ocean ![mondrian 1915](http://img.ahimsa.io/85rEC0DJiWJyTxOct2dxJI8od1yhcIb5WsYvxGiJ7pY=)","timestamp":1414193281},{"txid":"933c592a1a22b41a9a692aba57da649c91fe32403e8ff7b13f452071aa9820b9","board":"ahimsa-dev","author":"miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH","msg":"the mind is our medium","timestamp":1414017848,"blk":"00000000000016b6ff59b9fffcade68943bb02270b46d2a001054d95c56ca8ad","blkTimestamp":1414017952}]}`,
	},
	{
		endpoint: "/board/ahimsa-dev",
		body:     `{"summary":{"name":"ahimsa-dev","numBltns":4,"createdAt":1414017952,"lastActive":1414193281,"createdBy":"miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH"},"bltns":[{"txid":"f7800712c20377c2d29680c1aecf2331d6f80f5a44510d30ceb2e30fd5dafdcf","board":"ahimsa-dev","author":"mraY7GWs4G65ZYqWoPEqwPxb6aMuufTq2c","msg":"Here comes the sun","timestamp":1413079355},{"txid":"2963cc35727f4e2c2bd4186e4550fe82b204e446ff7096b425f236264e05c7c6","board":"ahimsa-dev","author":"miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH","msg":"pier and ocean ![mondrian 1915](http://img.ahimsa.io/85rEC0DJiWJyTxOct2dxJI8od1yhcIb5WsYvxGiJ7pY=)","timestamp":1414193281},{"txid":"933c592a1a22b41a9a692aba57da649c91fe32403e8ff7b13f452071aa9820b9","board":"ahimsa-dev","author":"miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH","msg":"the mind is our medium","timestamp":1414017848,"blk":"00000000000016b6ff59b9fffcade68943bb02270b46d2a001054d95c56ca8ad","blkTimestamp":1414017952},{"txid":"b0a1ba6e40d8f35aac526eecbc05d82b2a6d3c8d6a316627f593cbe592a777be","board":"ahimsa-dev","author":"mnPZBNTrLoCoSkAgSfKeeCujU3129PG6vn","msg":"","timestamp":1413216499,"blk":"00000000777213b4fd7c5d5a71b9b52608356c4194203b1b63d1bb0e6141d17d","blkTimestamp":1414813562,"bannedReason":"The Beatles are slanderous."}]}`,
	},
	// Test character encoding of board in bulletin, ensure that the output is utf-8 encoded
	{
		endpoint: "/bulletin/5ed76ba84d4116045df14ecf7a7eca86300a649ef3cbefdd2eeea3f84e1432dc",
		body:     `{"txid":"5ed76ba84d4116045df14ecf7a7eca86300a649ef3cbefdd2eeea3f84e1432dc","board":"#!~*Enc ded-boÄ\u0026\\/Ӂ","author":"mhDrE934aiWYESLKbxZjUsMBZBSHUbiZRw","msg":"Attempting to comply with RFC 3986. Россия","timestamp":1414897285}`,
	},
}

// Executes tests to verify that the json returned at an endpoint is correct
func TestResponses(t *testing.T) {

	ts := newTestServer(t)
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
			t.Logf("Responded with body:\n%s\nWanted:\n%s\n", body, testCase.body)
			t.Errorf("Bad json response")
		}
	}
}
