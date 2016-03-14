package lexer

import "encoding/json"

type MathGroup struct {
	Lbrace, Rbrace Leaf
	Nodes          []Node
}

func (s MathGroup) Pos() Pos {
	return s.Lbrace.Pos()
}

func (s MathGroup) End() Pos {
	return s.Rbrace.End()
}

func (s MathGroup) MarshalText() ([]byte, error) {
	var text []byte
	t, err := s.Lbrace.MarshalText()
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
	t, err = s.Rbrace.MarshalText()
	if err != nil {
		return nil, err
	}
	text = append(text, t...)
	return text, nil
}

func (s MathGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string
		Lbrace Node
		Nodes  []Node
		Rbrace Node
	}{
		"MathGroup",
		s.Lbrace,
		s.Nodes,
		s.Rbrace,
	})
}
