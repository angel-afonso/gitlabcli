package api

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func restReq(method string, path string, data []byte) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", rest, path), bytes.NewBuffer(data))

	if err != nil {
		log.Fatal(err)
	}

	return req
}

// Post send post request to gitlan api v4
// and bind the response in the given bind parameter
func (c *Client) Post(path string, data []byte, bind interface{}) {
	bindRestResponse(c.send(restReq("POST", path, data)), bind)
}

// Get send get request to gitlab api v4
// and bind the response in the given bind parameter
func (c *Client) Get(path string, bind interface{}) {
	bindRestResponse(c.send(restReq("GET", path, nil)), bind)
}
