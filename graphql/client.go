package graphql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/angel-afonso/gitlabcli/auth"
)

const (
	url = "https://gitlab.com/api/graphql"
)

// Client graphql client
type Client struct {
	session *auth.Session
}

type wrapper struct {
	Data interface{}
}

// NewClient create new graphql client
func NewClient(session *auth.Session) Client {
	return Client{session}
}

func (c *Client) send(data *strings.Reader, bind interface{}) {
	req, err := http.NewRequest("POST", url, data)

	if err != nil {
		log.Fatal(err)
	}

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

	response := wrapper{Data: bind}
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatal(err)
	}
}
