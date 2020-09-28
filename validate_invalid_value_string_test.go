package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateYamlInvalidValueString(t *testing.T) {
	t.Run("should not fail on usage of allowed values strings", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"bar":      "map[int]string",
			"baz": map[interface{}]interface{}{
				"ban": "[]int32",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of special characters in values", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "in+t",
			"bar":      "map[int]st&ring",
			"baz": map[interface{}]interface{}{
				"ban": "[]in@t32",
				"fan": "@",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorInvalidValueString("in+t", "foo", "root"),
			newValidationErrorInvalidValueString("map[int]st&ring", "bar", "root"),
			newValidationErrorInvalidValueString("[]in@t32", "ban", "baz"),
			newValidationErrorInvalidValueString("@", "fan", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of spaces characters in values", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "in t",
			"baz": map[interface{}]interface{}{
				"ban": "[]in t32",
				"fan": " ",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorInvalidValueString("in t", "foo", "root"),
			newValidationErrorInvalidValueString("[]in t32", "ban", "baz"),
			newValidationErrorInvalidValueString(" ", "fan", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of '*' characters in the wrong places", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"a":        "*string",
			"b":        "map[int]*string",
			"foo":      "int*",
			"bar":      "map[int*]string",
			"baz": map[interface{}]interface{}{
				"ban": "[*]int32",
				"fan": "*",
				"c":   "map[int]string*",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorInvalidValueString("int*", "foo", "root"),
			newValidationErrorInvalidValueString("map[int*]string", "bar", "root"),
			newValidationErrorInvalidValueString("[*]int32", "ban", "baz"),
			newValidationErrorInvalidValueString("*", "fan", "baz"),
			newValidationErrorInvalidValueString("map[int]string*", "c", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of '['/']' in the wrong places", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"a":        "[]string",
			"b":        "map[int]string]",
			"foo":      "int[]",
			"bar":      "[]map[int]string",
			"baz": map[interface{}]interface{}{
				"ban": "[]in[t32",
				"fan": "[]",
				"c":   "map[int][]string",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorInvalidValueString("map[int]string]", "b", "root"),
			newValidationErrorInvalidValueString("int[]", "foo", "root"),
			newValidationErrorInvalidValueString("[]in[t32", "ban", "baz"),
			newValidationErrorInvalidValueString("[]", "fan", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})
}
