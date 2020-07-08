package api

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFields(t *testing.T) {
	var query struct {
		Projects struct {
			Nodes []struct {
				ID          string `graphql-bind:"id"`
				Name        string
				Description string
			}
		} `graphql:"(membership: true)"`
	}

	structType := reflect.TypeOf(query)

	q := "{"

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		q += parseField(field)
	}
	q += "}"

	assert.Equal(t, "{projects(membership: true){nodes{id,name,description,}}}", q)

}
func TestFormatQuery(t *testing.T) {
	var query struct {
		Projects struct {
			Nodes []struct {
				Name string
			}
		} `graphql:"(membership: true)"`
	}

	q := formatQuery(query, nil)

	assert.Equal(t, `{"query":"{projects(membership: true){nodes{name,}}}","variables":{}}`, q)
}

func TestFormatQueryWithComposition(t *testing.T) {
	type composition struct {
		Name     string
		Lastname string
	}

	var query struct {
		Projects struct {
			Nodes []struct {
				composition   `graphql:"inner"`
				NoComposition int
			}
		} `graphql:"(membership: true)"`
	}

	q := formatQuery(query, nil)

	assert.Equal(t, `{"query":"{projects(membership: true){nodes{name,lastname,noComposition,}}}","variables":{}}`, q)
}

func TestFormatMutation(t *testing.T) {
	var query struct {
		MergeRequestCreate struct {
			MergeRequest struct {
				Title string
			}
		} `graphql:"(title:$title,projectPath:$path)"`
	}
	vars := struct {
		Title string `graphql-type:"String!"`
		Path  string `graphql-type:"String!"`
	}{
		Title: "asd",
		Path:  "asd",
	}
	q := formatMutation(query, vars)

	assert.Equal(t, `{"query":"mutation($title:String!,$path:String!,){mergeRequestCreate(title:$title,projectPath:$path){mergeRequest{title,}}}","variables":{"title":"asd","path":"asd"}}`, q)
}

func TestFormatMutationWithArrayVar(t *testing.T) {
	var query struct {
		MergeRequestAssing struct {
			MergeRequest struct {
				Title string
			}
		} `graphql:"(title:$title,usernames:$usernames)"`
	}
	vars := struct {
		Title     string   `graphql-type:"String!"`
		Usernames []string `graphql-type:"[String!]!"`
	}{
		Title:     "asd",
		Usernames: []string{`"asd"`},
	}
	q := formatMutation(query, vars)

	assert.Equal(t, `{"query":"mutation($title:String!,$usernames:[String!]!,){mergeRequestAssing(title:$title,usernames:$usernames){mergeRequest{title,}}}","variables":{"title":"asd","usernames":["asd"]}}`, q)
}

func TestQueryWithVariables(t *testing.T) {
	var query struct {
		Projects struct {
			Nodes []struct {
				Name string
			}
		} `graphql:"(membership: true)"`
	}

	q := formatQuery(query, struct {
		Field int32
		Foo   string
		Baz   int `graphql-type:"ID!"`
	}{
		Field: 123,
		Foo:   "asd",
		Baz:   123,
	})

	assert.Equal(t, `{"query":"query($field:Int,$foo:String,$baz:ID!,){projects(membership: true){nodes{name,}}}","variables":{"field":123,"foo":"asd","baz":123}}`, q)
}
