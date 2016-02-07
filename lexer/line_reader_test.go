package lexer_test

import (
	"io"
	"testing"

	"github.com/dmage/sherry/lexer"
)

func assertReadLine(t *testing.T, reader lexer.LineReader, want string) {
	got, err := reader.ReadLine()
	if err != nil {
		t.Errorf("got error %v", err)
		return
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertReadEOF(t *testing.T, reader lexer.LineReader) {
	_, err := reader.ReadLine()
	if err != io.EOF {
		t.Errorf("got %v, want %v", err, io.EOF)
	}
}

func TestStringLineReader(t *testing.T) {
	l := &lexer.StringLineReader{
		Input: "hello world",
	}
	assertReadLine(t, l, "hello world")
	assertReadEOF(t, l)

	l = &lexer.StringLineReader{
		Input: "a line\nanother line\n",
	}
	assertReadLine(t, l, "a line\n")
	assertReadLine(t, l, "another line\n")
	assertReadEOF(t, l)

	l = &lexer.StringLineReader{
		Input: "foobarbaz\nwidowed line",
	}
	assertReadLine(t, l, "foobarbaz\n")
	assertReadLine(t, l, "widowed line")
	assertReadEOF(t, l)
}
