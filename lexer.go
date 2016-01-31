package main

import (
	"log"

	"github.com/dmage/sherry/lexer"
)

func main() {
	l := lexer.Lexer{
		Input: []byte("echo|>2 cat"),
	}
	for {
		p, err := l.Get()
		if err != nil {
			log.Fatal(err)
		}
		if p == nil {
			break
		}
		log.Printf("%+v", p)
	}
}
