package lexer

import "bytes"

// Lexer is a splitter of the shell syntax.
type Lexer struct {
	Input    []byte
	consumed int
}

func (l *Lexer) consumeString(s string, kind Kind) *Leaf {
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

func (l *Lexer) consumeFunc(f func(c byte) bool, kind Kind) *Leaf {
	buf := l.Input[l.consumed:]
	i := 0
	for i < len(buf) && f(buf[i]) {
		i++
	}
	if i == 0 {
		return nil
	}
	leaf := &Leaf{
		Kind: kind,
		Data: buf[0:i],
		pos:  Pos(l.consumed),
	}
	l.consumed += i
	return leaf
}

func (l *Lexer) consumeWhile(b []byte, kind Kind) *Leaf {
	return l.consumeFunc(func(c byte) bool {
		return bytes.IndexByte(b, c) != -1
	}, kind)
}

func (l *Lexer) consumeUntil(b []byte, kind Kind) *Leaf {
	return l.consumeFunc(func(c byte) bool {
		return bytes.IndexByte(b, c) == -1
	}, kind)
}

// Get returns Node from unconsumed Input.
func (l *Lexer) Get() (Node, error) {
	if l.consumed >= len(l.Input) {
		return nil, nil
	}

	switch l.Input[l.consumed] {
	case ' ', '\t':
		leaf := l.consumeWhile([]byte(" \t"), Space)
		if leaf == nil {
			panic("unexpected nil")
		}
		return leaf, nil
	case '>':
		leaf := l.consumeString(">|", Operator)
		if leaf != nil {
			return leaf, nil
		}
		fallthrough
	case '|':
		l.consumed++
		return &Leaf{
			Kind: Operator,
			Data: l.Input[l.consumed-1 : l.consumed],
			pos:  Pos(l.consumed - 1),
		}, nil
	}

	leaf := l.consumeUntil([]byte("|> "), Word)
	if leaf != nil {
		return leaf, nil
	}

	panic("unreachable")
}
