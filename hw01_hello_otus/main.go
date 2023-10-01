package main

import (
	"fmt"
	"golang.org/x/example/stringutil"
)

var phrase = "Hello, OTUS!"

func main() {
	fmt.Println(stringutil.Reverse(phrase))
}
