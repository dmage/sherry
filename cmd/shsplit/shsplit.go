package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/dmage/sherry/lexer"
)

func main() {
	input := os.Stdin
	if len(os.Args) == 2 {
		s, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer input.Close()
		input = s
	} else if len(os.Args) != 1 {
		log.Fatalf("Usage: %s [<filename>]", os.Args[0])
	}

	buf, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatal(err)
	}

	l := lexer.Lexer{
		Input: buf,
	}
	output := []string{}
	for {
		p, err := l.Get()
		if err != nil {
			log.Fatal(err)
		}
		if p == nil {
			break
		}
		text, err := p.MarshalText()
		if err != nil {
			log.Fatal(err)
		}
		output = append(output, string(text))
	}
	err = json.NewEncoder(os.Stdout).Encode(output)
	if err != nil {
		log.Fatal(err)
	}
}
