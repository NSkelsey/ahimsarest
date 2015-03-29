package ahimsarest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/soapboxsys/ombudslib/pubrecdb"

	"testing"
)

func newTestServer(t *testing.T) *httptest.Server {
	db, err := pubrecdb.SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	handler := Handler("/", db)
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
	{"/author/0000000000000000000000000000000000", 204},
	{"/board/this-One-Isnt-Real", 404},
	// Ensure that a utf-8 url-encoded board is reachable
	{"/board/%23%21~%2AEnc%20ded-bo%C3%84&%5C/%D3%81", 200},
	{"/nilboard", 200},
	{"/noboard", 404},
	{"/recent", 200},
	{"/boards", 200},
	{"/unconfirmed", 200},
	{"/blocks/02-01-2006", 404},
	{"/blocks/01-11-2014", 200},
	{"/blocks/111-990-2014", 404},
	{"/status", 200},
	{"/authors", 200},
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
			t.Logf("In test: %s\n", testCase.endpoint)
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
	{
		endpoint: "/recent",
		body:     `[{"txid":"5df96dcb607701d19f7ae3a5da2708d834df7dc8ff505d74aa27dc82aeb7b3c1","board":"recent-test","author":"n1j3AYj82gnWmLnmFbTcF4GDxHNWNGyxG1","msg":"This is a test to see if recent confirmations works in the expected way.","timestamp":1415854832,"blk":"00000000459c5532eefae667ca4443c813c6c9a054cb2775ba14fd7a002890ff","blkTimestamp":1415862580}]`,
	},
	{
		endpoint: "/unconfirmed",
		body:     `[{"txid":"f7800712c20377c2d29680c1aecf2331d6f80f5a44510d30ceb2e30fd5dafdcf","board":"ahimsa-dev","author":"mraY7GWs4G65ZYqWoPEqwPxb6aMuufTq2c","msg":"Here comes the sun","timestamp":1413079355},{"txid":"2963cc35727f4e2c2bd4186e4550fe82b204e446ff7096b425f236264e05c7c6","board":"ahimsa-dev","author":"miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH","msg":"pier and ocean ![mondrian 1915](http://img.ahimsa.io/85rEC0DJiWJyTxOct2dxJI8od1yhcIb5WsYvxGiJ7pY=)","timestamp":1414193281},{"txid":"5ed76ba84d4116045df14ecf7a7eca86300a649ef3cbefdd2eeea3f84e1432dc","board":"#!~*Enc ded-boÄ\u0026\\/Ӂ","author":"mhDrE934aiWYESLKbxZjUsMBZBSHUbiZRw","msg":"Attempting to comply with RFC 3986. Россия","timestamp":1414897285},{"txid":"126484de57d01ab12ae19dfc7c4eb74087e6abb8e749badecc75d570ad577fa3","author":"mxmvvxMNaXvPPnU5vHXPoPEsrHbbnSAehh","msg":"This should be in the nil board.","timestamp":1414900834}]`,
	},
	{
		endpoint: "/nilboard",
		body:     `{"summary":{"name":"","numBltns":1,"createdAt":0,"lastActive":1414900834,"createdBy":"mxmvvxMNaXvPPnU5vHXPoPEsrHbbnSAehh"},"bltns":[{"txid":"126484de57d01ab12ae19dfc7c4eb74087e6abb8e749badecc75d570ad577fa3","author":"mxmvvxMNaXvPPnU5vHXPoPEsrHbbnSAehh","msg":"This should be in the nil board.","timestamp":1414900834}]}`,
	},
	{
		endpoint: "/blocks/01-11-2014",
		body:     `[{"hash":"000000009eca8c144e1be1daddee437e14d92c2379ed70adbb586c6b9a5610f4","prevHash":"0000000083f4f28cd2061073754383baeb38d3639e153b18e643268771e27f12","timestamp":1414801097,"height":305694,"numBltns":0},{"hash":"00000000efaee711979fe42e667188e50b1096e4d9cfcbc9a82101336189c2ca","prevHash":"00000000ef99c1e689c70bf2eaddbef5dc41412dfc0c350226d9caa850da307c","timestamp":1414800258,"height":305698,"numBltns":0},{"hash":"000000002f21b1943beb5c07a35fb89238b6dcd42312d39789b4d1b19b83f08a","prevHash":"00000000efaee711979fe42e667188e50b1096e4d9cfcbc9a82101336189c2ca","timestamp":1414801459,"height":305699,"numBltns":0},{"hash":"00000000777213b4fd7c5d5a71b9b52608356c4194203b1b63d1bb0e6141d17d","prevHash":"00000000f1cb8c224acdb0e1becbfa4218f1e13b0d4dbbce64d0a3c15d8bf55f","timestamp":1414813562,"height":305724,"numBltns":1}]`,
	},
	{
		endpoint: "/authors",
		body:     `[{"addr":"mhDrE934aiWYESLKbxZjUsMBZBSHUbiZRw","numBltns":1},{"addr":"miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH","numBltns":2,"firstBlkTs":1414017952},{"addr":"mnPZBNTrLoCoSkAgSfKeeCujU3129PG6vn","numBltns":1,"firstBlkTs":1414813562},{"addr":"mraY7GWs4G65ZYqWoPEqwPxb6aMuufTq2c","numBltns":1},{"addr":"mxmvvxMNaXvPPnU5vHXPoPEsrHbbnSAehh","numBltns":1},{"addr":"n1j3AYj82gnWmLnmFbTcF4GDxHNWNGyxG1","numBltns":1,"firstBlkTs":1415862580}]`,
	},
	{
		endpoint: "/block/00000000efaee711979fe42e667188e50b1096e4d9cfcbc9a82101336189c2ca",
		body:     `{"head":{"hash":"00000000efaee711979fe42e667188e50b1096e4d9cfcbc9a82101336189c2ca","prevHash":"00000000ef99c1e689c70bf2eaddbef5dc41412dfc0c350226d9caa850da307c","timestamp":1414800258,"height":305698,"numBltns":0},"bltns":[]}`,
	},
}

// Executes tests to verify that the json returned at an endpoint is correct
func TestResponses(t *testing.T) {

	ts := newTestServer(t)
	defer ts.Close()

	for _, testCase := range responseTests {
		url := ts.URL + testCase.endpoint
		res, err := http.Get(url)
		if err != nil {
			t.Logf("Endpoint: %s", testCase.endpoint)
			t.Logf("%v", res)
			t.Error(err)
		}
		if res.StatusCode != 200 {
			t.Logf("Endpoint: %s", testCase.endpoint)
			t.Errorf("Responded with a %d!", res.StatusCode)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Logf("Endpoint: %s", testCase.endpoint)
			t.Error(err)
		}
		if string(body) != testCase.body {
			t.Logf("Endpoint: %s", testCase.endpoint)
			t.Logf("Responded with body:\n%s\nWanted:\n%s\n", body, testCase.body)
			t.Errorf("Bad json response")
		}
	}
}
