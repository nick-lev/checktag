package main

import (
	"fmt"
	"reflect"
)

func checkTag(v interface{}) string {
	tagmap := make(map[string]int)
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fmt.Printf("!!!%v!!!\n", field)
			name := field.Name
			tag := string(field.Tag)
			if tag != "" {
				tagmap[tag]++
			}
			fmt.Println("\t->field:", field, "name:", name, "tag:", tag, "tagcount:", tagmap[tag])
			if tagmap[tag] > 1 {
				return fmt.Sprintf("duplicate tag: %s %s %d", tag, name, tagmap[tag])
			}
			if field.Type.Kind() == reflect.Struct {
				fmt.Println("internal struct finded!")
				if err := checkTag(field); err != "" {
					return err
				}
			}
		}
	}
	return ""
}

type tstruct struct {
	A  int
	B  string
	C  float64 `json:"page"`
	D  int     `json:"page2"`
	Dd struct {
		A int `json:"ppg1"`
		B int
		C int `json:"page"`
	} `json:"inernal struct page"`
	E int `json:"page3"`
}

func main() {
	i := tstruct{A: 10, B: "test"}
	fmt.Println(checkTag(i))
}
