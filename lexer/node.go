package lexer

// Pos is a position in the input data.
type Pos int

// Node is a high-level portion of the input data.
type Node interface {
	Pos() Pos
	End() Pos
	MarshalText() (text []byte, err error)
}
