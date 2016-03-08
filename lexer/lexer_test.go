package lexer_test

import (
	"testing"

	"github.com/dmage/sherry/lexer"
)

var lexerTests = []struct {
	Input  string
	Output []string
}{
	{
		"echo $foo|>2 cat",
		[]string{"echo", " ", "$foo", "|", ">", "2", " ", "cat"},
	},
	{
		"echo test >|file",
		[]string{"echo", " ", "test", " ", ">|", "file"},
	},
	{
		"echo \"$(date \"+%s\")\" >&2",
		[]string{"echo", " ", "\"$(date \"+%s\")\"", " ", ">&", "2"},
	},
	{
		"echo $(case $i in *) echo test; esac)",
		[]string{"echo", " ", "$(case $i in *) echo test; esac)"},
	},
	{
		"case $a in foo) case $b in *) echo foo$b;; esac;; esac",
		[]string{
			"case", " ", "$a", " ", "in", " ", "foo", ")", " ", "case", " ",
			"$b", " ", "in", " ", "*", ")", " ", "echo", " ", "foo", "$b",
			";;", " ", "esac", ";;", " ", "esac",
		},
	},
	{
		"echo case foo bar",
		[]string{
			"echo", " ", "case", " ", "foo", " ", "bar",
		},
	},
}

func TestLexer(t *testing.T) {
	for _, test := range lexerTests {
		l := lexer.Lexer{
			Input: []byte(test.Input),
		}
		pos := 0
		for _, lexeme := range test.Output {
			node, err := l.Get()
			if err != nil {
				t.Fatalf("%q: lexeme %q: %s", test.Input, lexeme, err)
			}

			if p := node.Pos(); int(p) != pos {
				t.Errorf("%q: lexeme %q: got position %d, await %d", test.Input, lexeme, p, pos)
			}

			text, err := node.MarshalText()
			if err != nil {
				t.Fatal(err)
			}
			if string(text) != lexeme {
				t.Errorf("%q: lexeme %q: got %q, await %q", test.Input, lexeme, text, lexeme)
			}
			pos += len(text)
		}
	}
}
