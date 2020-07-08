package api

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"gitlab.com/angel-afonso/gitlabcli/utils"
)

// graphqlReq generate request pointer
func graphqlReq(data *strings.Reader) *http.Request {
	req, err := http.NewRequest(post, graphql, data)

	if err != nil {
		log.Fatal(err)
	}

	return req
}

// Query Send a query graphql request
func (c *Client) Query(query interface{}, variables interface{}) {
	bindGraphqlResponse(c.send(graphqlReq(strings.NewReader(formatQuery(query, variables)))), query)
}

// Mutation Send a mutation graphql request
func (c *Client) Mutation(mutation interface{}, vars interface{}) {
	bindGraphqlResponse(c.send(graphqlReq(strings.NewReader(formatMutation(mutation, vars)))), mutation)
}

// generateQueryBody format a string as a graphql query body
func generateQueryBody(query interface{}) string {
	var structType reflect.Type

	if refl := reflect.TypeOf(query); refl.Kind() == reflect.Ptr {
		structType = refl.Elem()
	} else {
		structType = refl
	}

	return parseInnerFields(structType)
}

// formatMutation format a string as a graphql mutation
func formatMutation(mutation interface{}, vars interface{}) string {
	queryVars, variables := formatVariables(vars)

	return fmt.Sprintf(`{"query":"%s{%s}",%s}`,
		fmt.Sprintf("mutation(%s)", queryVars),
		generateQueryBody(mutation), variables)
}

// formatQuery returns a formated graphql query by stracting struct's fields
func formatQuery(query interface{}, vars interface{}) string {
	queryVars, variables := formatVariables(vars)

	return fmt.Sprintf(`{"query":"%s{%s}",%s}`,
		utils.Ternary(vars != nil, fmt.Sprintf("query(%s)", queryVars), ""),
		generateQueryBody(query), variables)
}

// parseField parse struct fields and generate the graphql query
func parseField(field reflect.StructField) string {
	q := ""

	switch tag := field.Tag.Get("graphql"); tag {
	case "inner":
		return parseInnerFields(field.Type)
	case "-":
		return q
	default:
		if bind, ok := field.Tag.Lookup("graphql-bind"); ok {
			q += bind
		} else {
			q += fmt.Sprintf("%s%s", string(bytes.ToLower([]byte{field.Name[0]})), field.Name[1:])
		}
		q += tag
	}

	switch field.Type.Kind() {
	case reflect.Slice:
		if field.Type.Elem().Kind() != reflect.Struct {
			q += ","
			break
		}
		q += fmt.Sprintf("{%s}", parseInnerFields(field.Type.Elem()))
		break
	case reflect.Ptr:
		q += fmt.Sprintf("{%s}", parseInnerFields(field.Type.Elem()))
		break
	case reflect.Struct:
		q += fmt.Sprintf("{%s}", parseInnerFields(field.Type))
		break

	default:
		q += ","
		break
	}

	return q
}

// parseInnerFields parse fields inside a strut field
func parseInnerFields(field reflect.Type) string {
	parsed := ""
	for i := 0; i < field.NumField(); i++ {
		parsed += parseField(field.Field(i))
	}
	return parsed
}

func formatVariables(vars interface{}) (string, string) {
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
			case reflect.Array, reflect.Slice:
				varType = fmt.Sprintf("[%s]", parseType(value.Kind()))
				varValue = fmt.Sprintf("%v", value)
				break
			case reflect.String:
				varType = "String,"
				varValue = fmt.Sprintf(`"%s"`, value)
				break
			default:
				varType = parseType(value.Kind())
				varValue = fmt.Sprintf("%v", value)
				break
			}

			if fieldType, ok := field.Tag.Lookup("graphql-type"); ok {
				varType = fieldType + ","
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

func parseArray(kind reflect.Kind) string {
	return fmt.Sprintf("[%s]", parseType(kind))
}

func parseType(kind reflect.Kind) string {
	switch kind {
	case reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
		return "Int,"
	case reflect.String:
		return "String,"
	case reflect.Bool:
		return "Boolean,"
	case reflect.Float32, reflect.Float64:
		return "Float,"
	}
	return ""
}
