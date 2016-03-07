package lexer

import "encoding/json"

type SubshellString struct {
	Lquote, Rquote Leaf
	Nodes          []Node
}

func (s SubshellString) Pos() Pos {
	return s.Lquote.Pos()
}

func (s SubshellString) End() Pos {
	return s.Rquote.End()
}

func (s SubshellString) MarshalText() ([]byte, error) {
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

func (s SubshellString) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string
		Lquote Node
		Nodes  []Node
		Rquote Node
	}{
		"SubshellString",
		s.Lquote,
		s.Nodes,
		s.Rquote,
	})
}
