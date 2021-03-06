package lexer

import "encoding/json"

type Word struct {
	Nodes []Node
}

func (s Word) Leaf(kind Kind) (Leaf, bool) {
	if len(s.Nodes) != 1 {
		return Leaf{}, false
	}
	if leaf, ok := s.Nodes[0].(Leaf); ok && leaf.Kind == kind {
		return leaf, true
	}
	return Leaf{}, false
}

func (s Word) Pos() Pos {
	return s.Nodes[0].Pos()
}

func (s Word) End() Pos {
	return s.Nodes[len(s.Nodes)-1].End()
}

func (s Word) MarshalText() ([]byte, error) {
	var text []byte
	for _, n := range s.Nodes {
		t, err := n.MarshalText()
		if err != nil {
			return nil, err
		}
		text = append(text, t...)
	}
	return text, nil
}

func (s Word) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string
		Nodes []Node
	}{
		"Word",
		s.Nodes,
	})
}
