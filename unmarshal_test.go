package yamltostruct

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	t.Run("should unmarshal without errors", func(t *testing.T) {
		yamlDataBytes := []byte(
			`_package: hello
foo: string
bar: 
  baz: int
  ban: "[]foo"`,
		)

		file, errs := Unmarshal(yamlDataBytes)

		assert.Equal(t, errs, []error{})

		output := normalizeWhitespace(printAST(file))
		expectedOutput := normalizeWhitespace(
			`package hello
			type bar struct {
				ban []foo
				baz int
			}
			type foo string
			`,
		)

		assert.Equal(t, output, expectedOutput)
	})

	t.Run("should return error", func(t *testing.T) {
		yamlDataBytes := []byte(
			`_package: hello
foo: string
bar: 
  baz: boo
  ban: "[]bool"`,
		)

		file, errs := Unmarshal(yamlDataBytes)
		assert.Equal(t, errs, []error{newValidationErrorTypeNotFound("boo", "bar")})
		var expectedFile *ast.File
		assert.Equal(t, file, expectedFile)
	})
}
