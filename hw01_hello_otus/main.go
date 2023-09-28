package main

import "fmt"

var phrase = "Hello, OTUS!"

func main() {
	var reverse []byte
	for i := len(phrase) - 1; i >= 0; i-- {
		reverse = append(reverse, phrase[i])
	}
	fmt.Println(string(reverse))
}
