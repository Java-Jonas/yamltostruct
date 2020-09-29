package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateYamlDataIllegalTypeName(t *testing.T) {
	t.Run("should not fail on valid key inputs", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"baz": map[interface{}]interface{}{
				"ban": "int",
			},
		}

		actualErrors := syntacticalValidation(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on spaces in key literal", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"fo o":     "int",
			"baz": map[interface{}]interface{}{
				"oof":  "int",
				"ba n": "int",
			},
		}

		actualErrors := syntacticalValidation(data)
		expectedErrors := []error{
			newValidationErrorIllegalTypeName("fo o", "root"),
			newValidationErrorIllegalTypeName("ba n", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should not fail on usage of allowed type names", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"bar":      "string",
			"baz": map[interface{}]interface{}{
				"ban": "int32",
			},
		}

		actualErrors := syntacticalValidation(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of keywords", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"break":    "int",
			"bar":      "string",
			"baz": map[interface{}]interface{}{
				"const": "int32",
			},
		}

		actualErrors := syntacticalValidation(data)
		expectedErrors := []error{
			newValidationErrorIllegalTypeName("break", "root"),
			newValidationErrorIllegalTypeName("const", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage special characters", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"*":        "int",
			"<":        "string",
			"fo$o":     "int",
			"baz": map[interface{}]interface{}{
				">-":    "int32",
				"bent{": "int32",
			},
		}

		actualErrors := syntacticalValidation(data)
		expectedErrors := []error{
			newValidationErrorIllegalTypeName("*", "root"),
			newValidationErrorIllegalTypeName("<", "root"),
			newValidationErrorIllegalTypeName("fo$o", "root"),
			newValidationErrorIllegalTypeName(">-", "baz"),
			newValidationErrorIllegalTypeName("bent{", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})
}
