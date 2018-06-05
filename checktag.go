package main

import (
	"fmt"
	"reflect"
	"regexp"
)

//is all duplication should be catched or we must stop at first dup?
func checkTag(v interface{}) string {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return ""
	}
	tagmap := make(map[string]int)
	var check func(t reflect.Type) string
	check = func(t reflect.Type) string {
		if t.Kind() != reflect.Struct {
			return ""
		}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			name := field.Name
			tag := string(field.Tag)
			if match, _ := regexp.MatchString("^json:", tag); !match { //is this case sensitive?
				tag = name
			}
			if match, _ := regexp.MatchString("^json:\"-\"$", tag); match {
				if field.Type.Kind() == reflect.Struct {
					err := check(field.Type)
					if err != "" {
						return err
					}
				}
				continue
			}
			re := regexp.MustCompile("^json:\"(.*)\"$")
			tag = re.ReplaceAllString(tag, "$1")

			re = regexp.MustCompile(",omitempty") //is space after coma possible?
			tag = re.ReplaceAllString(tag, "")
			if tag == "" {
				tag = name
			}
			tagmap[tag]++ //is case sensetivity for field name important for tag and struct field name?
			if tagmap[tag] > 1 {
				return fmt.Sprintf("duplicate tag:%s on field:%s", tag, name)
			}

			if field.Type.Kind() == reflect.Ptr {
				field.Type = field.Type.Elem()
			}
			if field.Type.Kind() == reflect.Struct {
				err := check(field.Type)
				if err != "" {
					return err
				}
			}
		}
		return ""
	}
	return check(t)
}

type First struct {
	A  int
	B  int `json:"b"`
	B2 int `json:"b,omitempty"` // conflict
	C  int `json:"-"`
	C2 int `json:"-,"`
	D  int `json:",omitempty"`
	E  int `json:"e,omitempty"`
	S  *Second
	T  *Third
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
