package graphql

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// Query Send a query graphql request
func (c *Client) Query(query interface{}, variables interface{}) {
	data := strings.NewReader(formatQuery(query, variables))
	c.send(data, query)
}

// formatQuery returns a formated graphql query
// by stracting struct's fields
func formatQuery(query interface{}, vars interface{}) string {
	q := `{"query":"`
	variables := ""

	var queryVars string
	queryVars, variables = parseVariables(vars)

	if vars != nil {
		q += fmt.Sprintf("query(%s)", queryVars)
	}

	q += "{"

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

	q += fmt.Sprintf(`}",%s}`, variables)
	return q
}

func parseVariables(vars interface{}) (string, string) {
	q := ""
	variables := `"variables":{`

	if vars != nil {

		structValue := reflect.ValueOf(vars)
		structType := reflect.TypeOf(vars)

		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			value := structValue.Field(i)

			name := fmt.Sprintf("%s%s", string(bytes.ToLower([]byte{field.Name[0]})), field.Name[1:])

			q += fmt.Sprintf("$%s:", name)
			variables += fmt.Sprintf(`"%s":`, name)

			var varType string
			var varValue string

			switch value.Kind() {
			case reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
				reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
				varType = "Int,"
				varValue = fmt.Sprintf("%d", value.Int())
				break
			case reflect.String:
				varType = "String,"
				varValue = fmt.Sprintf(`"%s"`, value.String())
				break
			case reflect.Bool:
				varType = "Boolean,"
				varValue = fmt.Sprintf("%t", value.Bool())
				break
			case reflect.Float32, reflect.Float64:
				varType = "Float,"
				varValue = fmt.Sprintf("%f", value.Float())
				break
			}

			if fieldType, ok := field.Tag.Lookup("graphql-type"); ok {
				varType = fieldType
			}

			if i != structType.NumField()-1 {
				varValue += ","
			}

			q += varType
			variables += varValue
		}
	}

	variables += "}"
	return q, variables
}

// parseFields parse struct fields and generate the graphql query
func parseFields(query reflect.StructField) string {
	q := ""
	if bind, ok := query.Tag.Lookup("graphql-bind"); ok {
		q += bind
	} else {
		q += fmt.Sprintf("%s%s", string(bytes.ToLower([]byte{query.Name[0]})), query.Name[1:])
	}

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
	default:
		q += ","
		break
	}

	return q
}
