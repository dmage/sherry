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
			"$b", " ", "in", " ", "*", ")", " ", "echo", " ", "foo$b", ";;",
			" ", "esac", ";;", " ", "esac",
		},
	},
	{
		"echo case foo bar",
		[]string{
			"echo", " ", "case", " ", "foo", " ", "bar",
		},
	},
	{
		"echo\n case foo in *) echo ;; esac",
		[]string{
			"echo", "\n", " ", "case", " ", "foo", " ", "in", " ", "*", ")",
			" ", "echo", " ", ";;", " ", "esac",
		},
	},
	{
		`"" case foo bar`,
		[]string{
			`""`, " ", "case", " ", "foo", " ", "bar",
		},
	},
	{
		"case $foo\"bar baz\" in *) echo ok; esac",
		[]string{"case", " ", "$foo\"bar baz\"", " ", "in", " ", "*", ")", " ", "echo", " ", "ok", ";", " ", "esac"},
	},
	{
		"case foo! in *) echo ok; esac",
		[]string{"case", " ", "foo!", " ", "in", " ", "*", ")", " ", "echo", " ", "ok", ";", " ", "esac"},
	},
	{
		">&2 case foo bar",
		// TODO(dmage): "case" is not a keyword here
		// TODO(dmage): ">&2" should be a single lexeme
		[]string{">&", "2", " ", "case", " ", "foo", " ", "bar"},
	},
	{
		"$(! case a in *) echo true; esac)",
		[]string{"$(! case a in *) echo true; esac)"},
	},
	{
		"$(!case a in *) echo true; esac)",
		[]string{"$(!case a in *)", " ", "echo", " ", "true", ";", " ", "esac", ")"},
	},
	{
		"{{ case foo bar }}",
		// TODO(dmage): "case" is not a keyword here
		[]string{"{{", " ", "case", " ", "foo", " ", "bar", " ", "}}"},
	},
	{
		"$({ case a in *) echo true; esac; })",
		[]string{"$({ case a in *) echo true; esac; })"},
	},
	{
		"echo foo\"b a r\"$baz",
		[]string{"echo", " ", "foo\"b a r\"$baz"},
	},
	{
		"echo \"f o o\\\"b a r\"",
		[]string{"echo", " ", "\"f o o\\\"b a r\""},
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
				t.Errorf("%q: lexeme %q: %s", test.Input, lexeme, err)
				break
			}
			if node == nil {
				t.Errorf("%q: got EOF, await %q", test.Input, lexeme)
				break
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

		if pos != len(test.Input) {
			t.Errorf("%q: unconsumed: %q", test.Input, test.Input[pos:])
		}
	}
}
