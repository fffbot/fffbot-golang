package main

import (
	"testing"
	"fmt"
)

func TestWithCode(t *testing.T) {
	url := "https://www.factorio.com/blog/post/fff-235"

	body, e := fetchPage(url)

	if e != nil {
		t.Error(e)
	}

	reply, e := parsePageToReply(url, body)

	if e != nil {
		t.Error(e)
	}

	fmt.Println(reply)
}
