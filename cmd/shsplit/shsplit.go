package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/dmage/sherry/lexer"
)

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
	nodes := []lexer.Node{}
	for {
		n, err := l.Get()
		if err != nil {
			log.Fatal(err)
		}
		if n == nil {
			break
		}
		nodes = append(nodes, n)
	}
	err = json.NewEncoder(os.Stdout).Encode(nodes)
	if err != nil {
		log.Fatal(err)
	}
}
