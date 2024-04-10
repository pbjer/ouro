package llm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

func NumberOfTokens(text string) (int, error) {
	tkm, err := tiktoken.EncodingForModel("gpt-4")
	if err != nil {
		err = fmt.Errorf("error getting encoding for model: %v", err)
		return 0, err
	}
	return len(tkm.Encode(text, nil, nil)), nil
}

func TrimNonJSON(s string) string {
	startIndex := strings.IndexAny(s, "{[")
	endIndex := strings.LastIndexAny(s, "]}")

	if startIndex == -1 || endIndex == -1 {
		return s // Return original string if no JSON boundaries are found
	}

	return s[startIndex : endIndex+1]
}

func StructToJSON(v interface{}) (string, error) {
	example := structToMap(reflect.TypeOf(v))
	jsonBytes, err := json.MarshalIndent(example, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func structToMap(t reflect.Type) interface{} {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		example := make(map[string]interface{})
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = field.Name
			}
			example[jsonTag] = structToMap(field.Type)
		}
		return example
	case reflect.Slice:
		return []interface{}{structToMap(t.Elem())}
	default:
		return getDefaultValue(t)
	}
}

func getDefaultValue(t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.Bool:
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return 0
	case reflect.Float32, reflect.Float64:
		return 0.0
	case reflect.String:
		return ""
	case reflect.Slice:
		return []interface{}{}
	default:
		return nil
	}
}
