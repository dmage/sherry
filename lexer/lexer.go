package lexer

import "bytes"

// Lexer is a splitter of the shell syntax.
type Lexer struct {
	Input    []byte
	consumed int
}

func (l *Lexer) consumeUntil(chars string, kind Kind) *Leaf {
	rest := l.Input[l.consumed:]
	i := bytes.IndexAny(rest, chars)
	if i == 0 {
		return nil
	}
	if i == -1 {
		i = len(rest)
	}
	leaf := &Leaf{
		Kind: kind,
		Data: rest[0:i],
		pos:  Pos(l.consumed),
	}
	l.consumed += i
	return leaf
}

func (l *Lexer) consume(s string, kind Kind) *Leaf {
	if l.consumed+len(s) > len(l.Input) {
		return nil
	}
	data := l.Input[l.consumed : l.consumed+len(s)]
	if string(data) != s {
		return nil
	}
	leaf := &Leaf{
		Kind: kind,
		Data: data,
		pos:  Pos(l.consumed),
	}
	l.consumed += len(data)
	return leaf
}

// Get returns Node from unconsumed Input.
func (l *Lexer) Get() (Node, error) {
	if l.consumed >= len(l.Input) {
		return nil, nil
	}

	switch l.Input[l.consumed] {
	case '>':
		leaf := l.consume(">|", Operator)
		if leaf != nil {
			return leaf, nil
		}
		fallthrough
	case '|', ' ':
		l.consumed++
		return &Leaf{
			Kind: Operator,
			Data: l.Input[l.consumed-1 : l.consumed],
			pos:  Pos(l.consumed - 1),
		}, nil
	}

	leaf := l.consumeUntil("|> ", Word)
	if leaf != nil {
		return leaf, nil
	}

	panic("unreachable")
}
