package domain

import (
	"fmt"
	"reflect"
	"strings"
)

func ParseSortParams(sort string) []SortObject {
	var sortParams []SortObject
	for _, part := range strings.Split(sort, ";") {
		fields := strings.Split(part, ",")
		if len(fields) == 2 {
			sortParams = append(sortParams, SortObject{
				SortBy:    fields[0],
				SortOrder: strings.ToUpper(fields[1]),
			})
		}
	}
	return sortParams
}

func GetOrderClause(sorts []SortObject) string {
	disallowedSort := map[string]bool{
		"item_name":           true,
		"size_us":             true,
		"condition":           true,
		"packaging_condition": true,
	}

	var orderClauses []string
	for _, sortObj := range sorts {
		if sortObj.SortBy == "" || (sortObj.SortOrder != "ASC" && sortObj.SortOrder != "DESC") {
			continue
		}
		if disallowedSort[sortObj.SortBy] {
			continue
		}
		orderClauses = append(orderClauses, sortObj.SortBy+" "+sortObj.SortOrder)
	}

	if len(orderClauses) > 0 {
		return " ORDER BY " + strings.Join(orderClauses, ", ")
	}
	return ""
}

func GetPaginationClause(page, pageSize int) string {
	if pageSize == 0 {
		pageSize = 10
	}
	offset := page * pageSize
	return fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
}

func MapToKeyValueArrays(req interface{}) ([]string, []interface{}) {
	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	keys := []string{}
	values := []interface{}{}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		val := v.Field(i)

		if !val.CanInterface() {
			continue
		}

		// Skip string kosong
		if val.Kind() == reflect.String && val.String() == "" {
			continue
		}

		// Skip pointer nil
		if val.Kind() == reflect.Ptr && val.IsNil() {
			continue
		}

		keys = append(keys, field.Name)
		values = append(values, v.Field(i).Interface())
	}
	return keys, values
}
