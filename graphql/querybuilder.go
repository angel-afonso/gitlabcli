package graphql

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// Query Send a query graphql request
func (c *Client) Query(query interface{}) {
	data := strings.NewReader(formatQuery(query))
	c.send(data, query)
}

// formatQuery returns a formated graphql query
// by stracting struct's fields
func formatQuery(query interface{}) string {
	q := `{"query":"{`

	var structType reflect.Type

	if reflect.TypeOf(query).Kind() == reflect.Ptr {
		structType = reflect.TypeOf(query).Elem()

	} else {
		structType = reflect.TypeOf(query)
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		q += parseFields(field)
	}

	q += `}"}`
	return q
}

// parseFields parse struct fields and generate the graphql query
func parseFields(query reflect.StructField) string {
	q := ""
	q += fmt.Sprintf("%s%s", string(bytes.ToLower([]byte{query.Name[0]})), query.Name[1:])

	if params, ok := query.Tag.Lookup("graphql"); ok {
		q += params
	}

	switch query.Type.Kind() {
	case reflect.Slice:
		q += "{"
		for i := 0; i < query.Type.Elem().NumField(); i++ {
			field := query.Type.Elem().Field(i)
			q += parseFields(field)
		}
		q += "}"
		break
	case reflect.Struct:
		q += "{"
		for i := 0; i < query.Type.NumField(); i++ {
			field := query.Type.Field(i)
			q += parseFields(field)
		}
		q += "}"
		break

	}

	return q
}
