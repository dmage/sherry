package lexer

import "encoding/json"

//go:generate stringer -type=Kind

// Kind allows to distinguish between different types of leaves.
type Kind int

// Leaf kinds.
const (
	Unknown Kind = iota
	Word
	Operator
	Space
	NewLine
	Comment
	Variable
	Quote
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

func (l Leaf) MarshalText() ([]byte, error) {
	return l.Data, nil
}

func (l Leaf) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
		Kind string
		Data string
	}{
		"Leaf",
		l.Kind.String(),
		string(l.Data),
	})
}
