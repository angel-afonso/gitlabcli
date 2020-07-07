package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gitlab.com/angel-afonso/gitlabcli/auth"
)

const (
	url     = "https://gitlab.com/api/graphql"
	graphql = "https://gitlab.com/api/graphql"
	rest    = "https://gitlab.com/api/v4"
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

func bindGraphqlResponse(body []byte, bind interface{}) {
	response := wrapper{Data: bind}
	err := json.Unmarshal(body, &response)

	if err != nil {
		log.Fatal(err)
	}

	if len(response.Errors) > 0 {
		for _, err := range response.Errors {
			log.Fatal(err.Message)
		}
	}
}

func bindRestResponse(body []byte, bind interface{}) {
	response := bind
	err := json.Unmarshal(body, &response)

	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) send(req *http.Request) []byte {
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", c.session.Type, c.session.Token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return body
}
