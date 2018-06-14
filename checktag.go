package main

import (
	"fmt"
	"reflect"
	"strings"
)

func tag2key(tag, path, name string) (key string) {
	if tag == "" {
		key = path + "." + name
	} else {
		s := strings.Split(tag, ",")
		switch {
		case s[0] == "": //fieldname as tag
			key = path + "." + name
		case s[0] == "-" && len(s) > 1: //- as tag accepted if content:"-,"
			key = path + "." + s[0]
		case s[0] == "-" && len(s) == 1: //ignore field if tag "-"
			key = ""
		default:
			key = path + "." + s[0]
		}
	}
	return key
}

func check(t reflect.Type, tagmap map[string]int, path string) {
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		tag := field.Tag.Get("json")
		key := tag2key(tag, path, name)
		if key == "" { // ignore empty
			continue
		}
		(tagmap)[key]++
		if field.Type.Kind() == reflect.Ptr {
			field.Type = field.Type.Elem()
		}
		if field.Type.Kind() == reflect.Struct {
			if field.Anonymous {
				key = path
			}
			check(field.Type, tagmap, key)
		}
	}
	return
}

func CheckTag(v interface{}) error {
	t := reflect.TypeOf(v)
	var err error
	if t.Kind() != reflect.Struct {
		return err
	}
	tagmap := make(map[string]int)
	path := "" // label for root path
	check(t, tagmap, path)
	for tag, val := range tagmap {
		if val > 1 {
			if err == nil {
				err = fmt.Errorf("%s", tag)
			} else {
				err = fmt.Errorf("%s,%s", err, tag)
			}
		}
	}
	return err
}
