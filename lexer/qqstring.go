package lexer

import "encoding/json"

// Leaf is a basic node type. Represents a piece of the input data.
type QQString struct {
	Lquote, Rquote Leaf
	Nodes          []Node
}

func (s QQString) Pos() Pos {
	return s.Lquote.Pos()
}

func (s QQString) End() Pos {
	return s.Rquote.End()
}

func (s QQString) MarshalText() ([]byte, error) {
	var text []byte
	t, err := s.Lquote.MarshalText()
	if err != nil {
		return nil, err
	}
	text = append(text, t...)
	for _, n := range s.Nodes {
		t, err := n.MarshalText()
		if err != nil {
			return nil, err
		}
		text = append(text, t...)
	}
	t, err = s.Rquote.MarshalText()
	if err != nil {
		return nil, err
	}
	text = append(text, t...)
	return text, nil
}

func (s QQString) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string
		Lquote Node
		Nodes  []Node
		Rquote Node
	}{
		"QQString",
		s.Lquote,
		s.Nodes,
		s.Rquote,
	})
}
