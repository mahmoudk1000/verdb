package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func Format(v any) (string, error) {
	var sb strings.Builder

	val := reflect.ValueOf(v)
	var sliceVal reflect.Value
	var elemType reflect.Type

	switch val.Kind() {
	case reflect.Struct:
		sliceVal = reflect.MakeSlice(reflect.SliceOf(val.Type()), 1, 1)
		sliceVal.Index(0).Set(val)
		elemType = val.Type()
	case reflect.Slice:
		if val.Len() == 0 {
			return "<empty>", nil
		}
		elemType = val.Index(0).Type()
		if elemType.Kind() != reflect.Struct {
			return "", fmt.Errorf("provided argument is not a struct or slice of structs")
		}
		sliceVal = val
	default:
		return "", fmt.Errorf("provided argument is not a struct or slice of structs")
	}

	numFields := elemType.NumField()
	headers := make([]string, numFields)
	colWidths := make([]int, numFields)

	for i := range numFields {
		headers[i] = elemType.Field(i).Name
		colWidths[i] = len(headers[i])
	}

	rows := make([][]string, sliceVal.Len())
	for r := 0; r < sliceVal.Len(); r++ {
		row := make([]string, numFields)
		v := sliceVal.Index(r)
		for c := range numFields {
			cell := fmt.Sprintf("%v", v.Field(c).Interface())
			row[c] = cell
			if len(cell) > colWidths[c] {
				colWidths[c] = len(cell)
			}
		}
		rows[r] = row
	}

	for i, h := range headers {
		sb.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], strings.ToUpper(h)))
	}
	sb.WriteString("\n")

	for _, row := range rows {
		for i, cell := range row {
			sb.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], cell))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func FormatJSON(v any) (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func FormatYAML(v any) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
