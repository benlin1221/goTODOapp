package utils

import (
	"reflect"
	"strings"
	"unicode"
)

// returns the first FieldName which has a matching json tag on
// model. if no match is found returns nil.
func FieldForJsonTag(tag string, model interface{}) *string {
	if len(tag) == 0 {
		return nil
	}
	name := ParseTag(tag)
	if !isValidTag(name) {
		name = ""
	}
	val := reflect.ValueOf(model)
	for i := 0; i < val.Type().NumField(); i++ {
		if val.Type().Field(i).Tag.Get("json") == name {
			fieldName := val.Type().Field(i).Name
			return &fieldName
		}
	}
	return nil
}

// parseTag returns a struct field's json tag value
// ommiting its comma-separated options.
func ParseTag(tag string) string {
	tag, _, _ = Cut(tag, ",")
	return tag
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}
