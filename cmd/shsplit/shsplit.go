package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/dmage/sherry/lexer"
)

type encoder struct {
	writer io.Writer
	depth  int
}

func (e *encoder) printPrefix(format string, args ...interface{}) {
	prefix := ""
	for i := 0; i < e.depth; i++ {
		prefix += ".  "
	}

	fmt.Fprintf(e.writer, prefix+format, args...)
}

func (e *encoder) printSuffix(format string, args ...interface{}) {
	fmt.Fprintf(e.writer, format+"\n", args...)
}

func (e *encoder) printf(format string, args ...interface{}) {
	e.printPrefix(format+"\n", args...)
}

var (
	byteSliceType = reflect.TypeOf([]byte{})
	leafType      = reflect.TypeOf(lexer.Leaf{})
)

func (e *encoder) encodeSlice(v reflect.Value) {
	e.printSuffix("[")
	e.depth++
	for i := 0; i < v.Len(); i++ {
		e.printPrefix("")
		e.Encode(v.Index(i).Interface())
	}
	e.depth--
	e.printf("]")
}

func (e *encoder) encodeStruct(v reflect.Value) {
	if v.Type() == leafType {
		leaf := v.Interface().(lexer.Leaf)
		e.printSuffix("Leaf %s { %q }", leaf.Kind, leaf.Data)
		return
	}

	e.printSuffix("%s {", v.Type().Name())
	e.depth++
	for i := 0; i < v.Type().NumField(); i++ {
		e.printPrefix("%s: ", v.Type().Field(i).Name)
		e.Encode(v.Field(i).Interface())
	}
	e.depth--
	e.printf("}")
}

func (e *encoder) Encode(i interface{}) {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Slice:
		e.encodeSlice(v)
	case reflect.Struct:
		e.encodeStruct(v)
	case reflect.Int:
		e.printf("%d", v.Int())
	case reflect.String:
		e.printf("%q", v.String())
	default:
		panic(v)
	}
}

var useJSON = flag.Bool("json", false, "dump as json")

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
	if *useJSON {
		err = json.NewEncoder(os.Stdout).Encode(nodes)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		e := &encoder{writer: os.Stdout}
		e.Encode(nodes)
	}
}
