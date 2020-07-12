package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gitlab.com/angel-afonso/gitlabcli/auth"
	"gopkg.in/gookit/color.v1"
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

func bindGraphqlResponse(body []byte, bind interface{}) {
	response := wrapper{Data: bind}
	err := json.Unmarshal(body, &response)

	if err != nil {
		color.Red.Printf("\n%s\n", err.Error())
		os.Exit(1)
	}

	if len(response.Errors) > 0 {
		for _, err := range response.Errors {
			color.Red.Println(err.Message)
			os.Exit(1)
		}
	}
}

func bindRestResponse(body []byte, bind interface{}) {
	response := bind
	err := json.Unmarshal(body, &response)

	if err != nil {
		color.Red.Println(err.Error())
		os.Exit(1)
	}
}

func (c *Client) send(req *http.Request) []byte {
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", c.session.Type, c.session.Token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		color.Red.Println(err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		color.Red.Println(err.Error())
		os.Exit(1)
	}

	return body
}
