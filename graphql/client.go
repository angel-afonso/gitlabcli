package graphql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", c.session.Type, c.session.Token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	response := wrapper{Data: bind}
	fmt.Printf(string(body))
	err = json.Unmarshal(body, &response)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", response.Data)
}
