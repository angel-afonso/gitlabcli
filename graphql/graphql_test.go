package graphql

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
		q += parseFields(field)
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

	assert.Equal(t, `{"query":"query($field:Int,$foo:String,$baz:ID!){projects(membership: true){nodes{name,}}}","variables":{"field":123,"foo":"asd","baz":123}}`, q)
}
