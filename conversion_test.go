package yamltostruct

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

// "ab c  de\nf" => "ab c de f"
func normalizeWhitespace(str string) string {
	var b strings.Builder
	b.Grow(len(str))

	var wroteSpace bool = true

	for _, ch := range str {
		var isSpace bool = unicode.IsSpace(ch)

		if isSpace && wroteSpace {
			continue
		}

		if isSpace {
			b.WriteRune(' ')
		} else {
			b.WriteRune(ch)
		}

		if isSpace {
			wroteSpace = true
		} else {
			wroteSpace = false
		}
	}

	return b.String()
}

func printDecls(decls []ast.Decl) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), decls)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func toNormalizedSourceCode(inputYamlData map[interface{}]interface{}, expectedSourceCode string) (string, string) {
	inputAST := convertToAST(inputYamlData)
	inputSourceCode := printDecls(inputAST.Decls)
	normalizedInput := normalizeWhitespace(inputSourceCode)
	normalizedExpectedOutput := normalizeWhitespace(expectedSourceCode)
	return normalizedInput, normalizedExpectedOutput
}

func TestConvertToASTBasicCases(t *testing.T) {
	t.Run("should convert named types", func(t *testing.T) {
		input := map[interface{}]interface{}{
			"foo": "string",
			"bar": "int",
		}
		expectedOutput := `
		type bar int
		type foo string`
		normalizedActualOutput, normalizedExpectedOutput := toNormalizedSourceCode(input, expectedOutput)

		assert.Equal(t, normalizedActualOutput, normalizedExpectedOutput)
	})
	t.Run("should convert struct types", func(t *testing.T) {
		input := map[interface{}]interface{}{
			"foo": map[interface{}]interface{}{
				"bar": "int",
			},
		}
		expectedOutput := `
		type foo struct{
			bar int
		}`
		normalizedActualOutput, normalizedExpectedOutput := toNormalizedSourceCode(input, expectedOutput)

		assert.Equal(t, normalizedActualOutput, normalizedExpectedOutput)
	})
}

func TestAlphabeticalRange(t *testing.T) {
	t.Run("should loop in alphabetical range", func(t *testing.T) {
		input := map[interface{}]interface{}{
			"a": "1",
			"b": "2",
			"c": "3",
		}

		for i := 0; i < 100; i++ {
			var receivedKeys []string
			var receivedValues []string
			alphabeticalRange(input, func(key string, value interface{}) {
				_key := fmt.Sprintf("%v", key)
				receivedKeys = append(receivedKeys, _key)
				_value := fmt.Sprintf("%v", value)
				receivedValues = append(receivedValues, _value)
			})
			assert.Equal(t, receivedKeys, []string{"a", "b", "c"})
			assert.Equal(t, receivedValues, []string{"1", "2", "3"})
		}

	})
}
