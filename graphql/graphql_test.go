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
				Name string
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

	assert.Equal(t, "{projects(membership: true){nodes{name}}}", q)

}
func TestFormatQuery(t *testing.T) {
	var query struct {
		Projects struct {
			Nodes []struct {
				Name string
			}
		} `graphql:"(membership: true)"`
	}

	q := formatQuery(query)

	assert.Equal(t, `{"query":"{projects(membership: true){nodes{name}}}"}`, q)
}
