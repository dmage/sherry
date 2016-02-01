package lexer

import "bytes"

// Lexer is a splitter of the shell syntax.
type Lexer struct {
	Input    []byte
	consumed int
}

func (l *Lexer) consume(size int, kind Kind) Leaf {
	leaf := Leaf{
		Kind: kind,
		Data: l.Input[l.consumed : l.consumed+size],
		pos:  Pos(l.consumed),
	}
	l.consumed += size
	return leaf
}

func (l *Lexer) tryConsumeString(s string, kind Kind) (Leaf, bool) {
	if l.consumed+len(s) > len(l.Input) {
		return Leaf{}, false
	}
	data := l.Input[l.consumed : l.consumed+len(s)]
	if string(data) != s {
		return Leaf{}, false
	}
	return l.consume(len(s), kind), true
}

func (l *Lexer) consumeFunc(f func(c byte) bool, kind Kind) Leaf {
	buf := l.Input[l.consumed:]
	i := 0
	for i < len(buf) && f(buf[i]) {
		i++
	}
	if i == 0 {
		panic("nothing consumed")
	}
	return l.consume(i, kind)
}

func (l *Lexer) consumeWhile(b []byte, kind Kind) Leaf {
	return l.consumeFunc(func(c byte) bool {
		return bytes.IndexByte(b, c) != -1
	}, kind)
}

func (l *Lexer) consumeUntil(b []byte, kind Kind) Leaf {
	return l.consumeFunc(func(c byte) bool {
		return bytes.IndexByte(b, c) == -1
	}, kind)
}

func isVariableName(c byte) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z'
}

func (l *Lexer) consumeVariable() Leaf {
	buf := l.Input[l.consumed:]
	if len(buf) == 0 || buf[0] != '$' {
		panic("consumeVariable expects $")
	}
	if len(buf) == 1 {
		return l.consume(1, Variable)
	}
	if !isVariableName(buf[1]) {
		return l.consume(2, Variable)
	}
	i := 1
	for i < len(buf) && isVariableName(buf[i]) {
		i++
	}
	return l.consume(i, Variable)
}

// Get returns Node from unconsumed Input.
func (l *Lexer) Get() (Node, error) {
	if l.consumed >= len(l.Input) {
		return nil, nil
	}

	switch l.Input[l.consumed] {
	case ' ', '\t':
		return l.consumeWhile([]byte(" \t"), Space), nil
	case '#':
		return l.consumeUntil([]byte("\n"), Comment), nil
	case '\n':
		return l.consume(1, NewLine), nil
	case '$':
		return l.consumeVariable(), nil
	case ';':
		if leaf, ok := l.tryConsumeString(";;", Operator); ok {
			return leaf, nil
		}
		return l.consume(1, Operator), nil
	case '&':
		if leaf, ok := l.tryConsumeString("&&", Operator); ok {
			return leaf, nil
		}
		return l.consume(1, Operator), nil
	case '|':
		if leaf, ok := l.tryConsumeString("||", Operator); ok {
			return leaf, nil
		}
		return l.consume(1, Operator), nil
	case '<':
		for _, op := range []string{"<<-", "<<", "<>", "<&"} {
			if leaf, ok := l.tryConsumeString(op, Operator); ok {
				return leaf, nil
			}
		}
		return l.consume(1, Operator), nil
	case '>':
		for _, op := range []string{">>", ">&", ">|"} {
			if leaf, ok := l.tryConsumeString(op, Operator); ok {
				return leaf, nil
			}
		}
		return l.consume(1, Operator), nil
	case '!', '(', ')', '{', '}':
		return l.consume(1, Operator), nil
	}

	return l.consumeUntil([]byte(" \t#\n$;&|<>!(){}"), Word), nil
}
