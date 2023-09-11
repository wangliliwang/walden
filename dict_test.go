package main

import (
	"fmt"
	"testing"
)

func TestECDict(t *testing.T) {
	dict := NewECDict()
	fmt.Println(dict.Match("me"))
	fmt.Println(dict.Match("like"))
}
