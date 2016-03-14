package lexer_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/dmage/sherry/lexer"
)

var mathTests = []struct {
	Input       string
	Output      []byte
}{
	{
		"$((2+1))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Variable","Data":"2"},{"Type":"Leaf","Kind":"Term","Data":"+"},{"Type":"Leaf","Kind":"Variable","Data":"1"}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
	{
		"$((2\\++))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Variable","Data":"2"},{"Type":"Leaf","Kind":"Escaped","Data":"\\+"},{"Type":"Leaf","Kind":"Term","Data":"+"}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
	{
		"$((123))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Variable","Data":"123"}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
	{
		"$((010))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Variable","Data":"010"}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
	{
		"$((0x80))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Variable","Data":"0x80"}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
	{
		"$((0x80x))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Variable","Data":"0x80"},{"Type":"Leaf","Kind":"Variable","Data":"x"}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
	{
		"$(( (2 << 1) + 3 ))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Variable","Data":"2"},{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"Leaf","Kind":"Term","Data":"\u003c\u003c"},{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"Leaf","Kind":"Variable","Data":"1"}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}},{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"Leaf","Kind":"Term","Data":"+"},{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"Leaf","Kind":"Variable","Data":"3"},{"Type":"Leaf","Kind":"Space","Data":" "}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
	{
		"$(( x=8 , y=12 ))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"Leaf","Kind":"Variable","Data":"x"},{"Type":"Leaf","Kind":"Term","Data":"="},{"Type":"Leaf","Kind":"Variable","Data":"8"},{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"Leaf","Kind":"Term","Data":","},{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"Leaf","Kind":"Variable","Data":"y"},{"Type":"Leaf","Kind":"Term","Data":"="},{"Type":"Leaf","Kind":"Variable","Data":"12"},{"Type":"Leaf","Kind":"Space","Data":" "}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
	{
		"$((a b))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Variable","Data":"a"},{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"Leaf","Kind":"Variable","Data":"b"}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
	{
		"$((+++i))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Term","Data":"++"},{"Type":"Leaf","Kind":"Term","Data":"+"},{"Type":"Leaf","Kind":"Variable","Data":"i"}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
	{
		"$((--x=7))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"MathGroup","Lbrace":{"Type":"Leaf","Kind":"Quote","Data":"("},"Nodes":[{"Type":"Leaf","Kind":"Term","Data":"--"},{"Type":"Leaf","Kind":"Variable","Data":"x"},{"Type":"Leaf","Kind":"Term","Data":"="},{"Type":"Leaf","Kind":"Variable","Data":"7"}],"Rbrace":{"Type":"Leaf","Kind":"Quote","Data":")"}}]}`),
	},
/*
	{
		"$((date +%s); (pwd))",
		[]byte(`{"Type":"Word","Nodes":[{"Type":"SubshellString","Lquote":{"Type":"Leaf","Kind":"Quote","Data":"$("},"Nodes":[{"Type":"Leaf","Kind":"Operator","Data":"("},{"Type":"Word","Nodes":[{"Type":"Leaf","Kind":"Term","Data":"date"}]},{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"Word","Nodes":[{"Type":"Leaf","Kind":"Term","Data":"+%s"}]}],"Rquote":{"Type":"Leaf","Kind":"Quote","Data":")"}}]},{"Type":"Leaf","Kind":"Operator","Data":";"},{"Type":"Leaf","Kind":"Space","Data":" "},{"Type":"Leaf","Kind":"Operator","Data":"("},{"Type":"Word","Nodes":[{"Type":"Leaf","Kind":"Term","Data":"pwd"}]},{"Type":"Leaf","Kind":"Operator","Data":")"},{"Type":"Leaf","Kind":"Operator","Data":")"}`),
	},
*/
}

func showDiff(t *testing.T, a, b string) {
	for i := 0; i < len(a); i += 1 {
		if i == len(b) {
			t.Log("unexpected tail:", a[i:])
			break
		}
		if a[i] != b[i] {
			t.Log("expect:", b[i:])
			t.Log("got   :", a[i:])
			break
		}
	}
}

func TestArithmetic(t *testing.T) {
	for _, test := range mathTests {
		l := lexer.Lexer{
			Input: []byte(test.Input),
		}

		node, err := l.Get()

		if err != nil {
			t.Errorf("%q: %v", test.Input, err)
			continue
		}

		if node == nil {
			t.Errorf("%q: got EOF", test.Input)
			continue
		}

		out, err := json.Marshal(node)
		if err != nil {
			t.Fatalf("%v", err)
		}

		if c := bytes.Compare(out, test.Output); c != 0 {
			t.Log("")
			t.Errorf("%q: differ", test.Input)

			strOut := string(out)
			strExpect := string(test.Output)

			if c >= 0 {
				showDiff(t, strOut, strExpect)
			} else {
				showDiff(t, strExpect, strOut)
			}
		}
	}
}
