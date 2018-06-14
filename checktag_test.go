package main

import (
	"fmt"
	"strings"
	"testing"
)

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

func TestCheckTag(t *testing.T) {
	var v First
	err := CheckTag(v)
	if len(strings.Split(fmt.Sprintf("%s", err), ",")) != 6 {
		t.Errorf("Wrong error count")
	}
}
