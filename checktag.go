package main

import (
	"fmt"
	"reflect"
	"strings"
)

func tag2key(tag, path, name string) string {
	if tag == "" {
		tag = path + "->" + name
	} else {
		s := strings.Split(tag, ",")
		switch {
		case s[0] == "": //fieldname as tag
			tag = path + "->" + name
		case s[0] == "-" && len(s) > 1: //- as tag accepted if content:"-,"
			tag = path + "->" + s[0]
		case s[0] == "-" && len(s) == 1: //ignore field if tag "-"
			tag = ""
		default:
			tag = path + "->" + s[0]
		}
	}
	return tag
}

func check(t reflect.Type, tagmap *map[string]int, path string) {
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		tag := field.Tag.Get("json")
		tag = tag2key(tag, path, name)
		if tag == "" { // ignore empty
			continue
		}
		(*tagmap)[tag]++
		if field.Type.Kind() == reflect.Ptr {
			field.Type = field.Type.Elem()
		}
		if field.Type.Kind() == reflect.Struct {
			if field.Anonymous {
				tag = path
			}
			check(field.Type, tagmap, tag)
		}
	}
	return

}

func checkTag(v interface{}) error {
	t := reflect.TypeOf(v)
	var err error
	if t.Kind() != reflect.Struct {
		return err
	}
	tagmap := make(map[string]int)
	path := "{Root}"
	check(t, &tagmap, path)
	for tag, val := range tagmap {
		if val > 1 {
			if err == nil {
				err = fmt.Errorf("%s", tag)
			} else {
				err = fmt.Errorf("%s, %s", err, tag)
			}
		}
	}
	return err
}

type First struct {
	A  int
	B  int `json:"b"`
	B2 int `json:"b,omitempty"` // conflict
	C  int `json:"-,"`
	C2 int `json:"-"`
	D  int `json:",omitempty"`
	E  int `json:"e,omitempty"`
	*Second
	Third *Third
}

type Second struct {
	A  int // conflict
	B3 int `json:"b"` // conflict
	C  int
	D  int //conflict
	E  int `json:"e"`     // conflict
	F  int `json:"Third"` // conflict
}

type Third struct {
	A int
	B int `json:"b"`
	C int `json:"F"`
	D int `json:"A"` // conflict
	E int `json:"B"`
	F int `json:"-"`
}

func main() {
	var v First
	fmt.Println(checkTag(v))
}
