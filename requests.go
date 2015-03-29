package ahimsarest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombjson"
)

type Client struct {
	httpCli *http.Client
	base    string
}

// NewClient creates an http client to communicate with a ombwebapp api
func NewClient(base string) *Client {
	return &Client{
		httpCli: &http.Client{},
		base:    base,
	}
}

func (c Client) GetJsonBlockHead(hash *wire.ShaHash) (ombjson.JsonBlkHead, error) {

	blkhead := ombjson.JsonBlkHead{}
	url := fmt.Sprintf("%s/api/blockhead/%s", c.base, hash.String())
	resp, err := c.httpCli.Get(url)
	if err != nil {
		return blkhead, err
	}

	defer resp.Body.Close()
	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return blkhead, err
	}

	if resp.StatusCode == 200 {
		if err = json.Unmarshal(blob, &blkhead); err != nil {
			return blkhead, err
		}
	} else {
		return blkhead, fmt.Errorf("Unexpected status code")
	}

	return blkhead, nil
}

func (c Client) GetJsonAuthor(addr btcutil.Address) (ombjson.AuthorResp, error) {

	authorResp := ombjson.AuthorResp{}
	url := fmt.Sprintf("%s/api/author/%s", c.base, addr.String())
	resp, err := c.httpCli.Get(url)
	if err != nil {
		return authorResp, err
	}

	defer resp.Body.Close()
	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return authorResp, err
	}

	// Unmarshal the json into an AuthorResp struct
	if resp.StatusCode == 200 {
		if err = json.Unmarshal(blob, &authorResp); err != nil {
			return authorResp, err
		}
	} else {
		if resp.StatusCode != 204 {
			err = fmt.Errorf("Server responded with: %s", resp.Status)
			return authorResp, err
		}
	}

	return authorResp, nil
}
