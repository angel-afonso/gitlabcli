package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/angel-afonso/gitlabcli/auth"
)

const (
	url     = "https://gitlab.com/api/graphql"
	graphql = "https://gitlab.com/api/graphql"
	rest    = "https://gitlab.com/api/v4"

	get  = "GET"
	post = "POST"
)

// Client graphql client
type Client struct {
	session *auth.Session
}

type wrapper struct {
	Data   interface{}
	Errors []struct {
		Message string
	}
}

// NewClient create new graphql client
func NewClient(session *auth.Session) Client {
	return Client{session}
}

func bindGraphqlResponse(body []byte, bind interface{}) error {
	response := wrapper{Data: bind}
	err := json.Unmarshal(body, &response)

	if err != nil {
		return err
	}

	if len(response.Errors) > 0 {
		for _, err := range response.Errors {
			return errors.New(err.Message)
		}
	}

	return nil
}

func bindRestResponse(body []byte, bind interface{}) error {
	response := bind
	err := json.Unmarshal(body, &response)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) send(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", c.session.Type, c.session.Token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
