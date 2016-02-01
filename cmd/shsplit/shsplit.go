package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/dmage/sherry/lexer"
)

var kind = flag.Bool("kind", false, "prefix lexemes by the kind")

func main() {
	flag.Parse()

	input := os.Stdin
	if len(flag.Args()) == 1 {
		s, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer input.Close()
		input = s
	} else if len(flag.Args()) != 0 {
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
		if *kind {
			if leaf, ok := p.(lexer.Leaf); ok {
				output = append(output, leaf.Kind.String()+":"+string(text))
			} else {
				output = append(output, "!Leaf:"+string(text))
			}
		} else {
			output = append(output, string(text))
		}
	}
	err = json.NewEncoder(os.Stdout).Encode(output)
	if err != nil {
		log.Fatal(err)
	}
}
