package lexer

import "fmt"

//go:generate stringer -type=Kind

// Kind allows to distinguish between different types of leaves.
type Kind int

// Leaf kinds.
const (
	Unknown Kind = iota
	Word
	Operator
	Space
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
	return fmt.Sprintf("{pos:%-6d Kind:%-10v Data:%q}", l.pos, l.Kind, l.Data)
}
