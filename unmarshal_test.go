package yamltostruct

import (
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
  ban: "[]bool"`,
		)

		file, errs := Unmarshal(yamlDataBytes)

		assert.Equal(t, errs, []error{})

		output := normalizeWhitespace(printAST(file))
		expectedOutput := normalizeWhitespace(
			`package hello
			type bar struct {
				ban []bool
				baz int
			}
			type foo string
			`,
		)

		assert.Equal(t, output, expectedOutput)
	})
}
