package api

import (
	"bytes"
	"fmt"
	"net/http"
)

func restReq(method string, path string, data []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", rest, path), bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	return req, nil
}

// Post send post request to gitlan api v4
// and bind the response in the given bind parameter
func (c *Client) Post(path string, data []byte, bind interface{}) error {
	req, err := restReq(post, path, data)

	if err != nil {
		return err
	}

	bytes, err := c.send(req)
	if err != nil {
		return err
	}

	bindRestResponse(bytes, bind)

	return nil
}

// Get send get request to gitlab api v4
// and bind the response in the given bind parameter
func (c *Client) Get(path string, bind interface{}) error {
	req, err := restReq(get, path, nil)
	if err != nil {
		return err
	}

	bytes, err := c.send(req)
	if err != nil {
		return err
	}

	bindRestResponse(bytes, bind)

	return nil
}
