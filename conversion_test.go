package yamltostruct

import (
	"bytes"
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

func printAST(ast *ast.File) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), ast)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func toNormalizedSourceCode(inputYamlData map[interface{}]interface{}, expectedSourceCode string) (string, string) {
	inputAST := convertToAST(inputYamlData)
	inputSourceCode := printAST(inputAST)
	normalizedInput := normalizeWhitespace(inputSourceCode)
	normalizedExpectedOutput := normalizeWhitespace(expectedSourceCode)
	return normalizedInput, normalizedExpectedOutput
}

func TestConvertToASTBasicCases(t *testing.T) {
	t.Run("should convert specified package name", func(t *testing.T) {
		input := map[interface{}]interface{}{
			"_package": "foobar",
		}
		expectedOutput := `
		package foobar
		`
		normalizedActualOutput, normalizedExpectedOutput := toNormalizedSourceCode(input, expectedOutput)

		assert.Equal(t, normalizedActualOutput, normalizedExpectedOutput)
	})
	t.Run("should convert named types", func(t *testing.T) {
		input := map[interface{}]interface{}{
			"_package": "foobar",
			"foo":      "string",
			"bar":      "int",
		}
		expectedOutput := `
		package foobar
		type foo string
		type bar int
		`
		normalizedActualOutput, normalizedExpectedOutput := toNormalizedSourceCode(input, expectedOutput)

		assert.Equal(t, normalizedActualOutput, normalizedExpectedOutput)
	})
	t.Run("should convert struct types", func(t *testing.T) {
		input := map[interface{}]interface{}{
			"_package": "foobar",
			"foo": map[interface{}]interface{}{
				"bar": "int",
			},
		}
		expectedOutput := `
		package foobar
		type foo struct{
			bar int
		}
		`
		normalizedActualOutput, normalizedExpectedOutput := toNormalizedSourceCode(input, expectedOutput)

		assert.Equal(t, normalizedActualOutput, normalizedExpectedOutput)
	})
}
