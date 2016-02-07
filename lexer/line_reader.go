package lexer

import (
	"io"
	"strings"
)

type LineReader interface {
	ReadLine() (string, error)
}

type StringLineReader struct {
	Input string
	pos   int
}

func (r *StringLineReader) ReadLine() (string, error) {
	if r.pos >= len(r.Input) {
		return "", io.EOF
	}

	buf := r.Input[r.pos:]
	n := strings.IndexByte(buf, '\n')
	if n == -1 {
		n = len(buf)
	} else {
		n += 1
	}
	r.pos += n
	return buf[:n], nil
}
