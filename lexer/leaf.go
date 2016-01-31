package lexer

import "fmt"

// Kind allows to distinguish between different types of leaves.
type Kind int

// Leaf kinds.
const (
	Unknown Kind = iota
	Word
	Operator
)

// Leaf is a basic node type. Represents a piece of the input data.
type Leaf struct {
	Kind Kind
	Data []byte
	pos  Pos
}

func (l Leaf) Pos() Pos {
	return l.pos
}

func (l Leaf) End() Pos {
	return l.pos + Pos(len(l.Data))
}

func (l Leaf) MarshalText() (text []byte, err error) {
	return l.Data, nil
}

func (l Leaf) String() string {
	return fmt.Sprintf("{Kind:%d Data:%q pos:%d}", l.Kind, l.Data, l.pos)
}
