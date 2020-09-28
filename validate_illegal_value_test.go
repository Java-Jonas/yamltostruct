package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateYamlInvalidValue(t *testing.T) {
	t.Run("should not fail on usage of allowed values", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"bar":      "string",
			"baz": map[interface{}]interface{}{
				"ban": "int32",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of numbers", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      1,
			"bar":      1.2,
			"baz": map[interface{}]interface{}{
				"ban": 3,
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorInvalidValue("foo", "root"),
			newValidationErrorInvalidValue("bar", "root"),
			newValidationErrorInvalidValue("ban", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of empty and nil values", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      nil,
			"bar":      "",
			"baz": map[interface{}]interface{}{
				"ban": nil,
				"baf": "",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorInvalidValue("foo", "root"),
			newValidationErrorInvalidValue("bar", "root"),
			newValidationErrorInvalidValue("ban", "baz"),
			newValidationErrorInvalidValue("baf", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of invalid list values", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"bar":      "string",
			"baz": map[interface{}]interface{}{
				"ban":  "int32",
				"mant": []interface{}{},
			},
			"rant": []interface{}{},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorInvalidValue("mant", "baz"),
			newValidationErrorInvalidValue("rant", "root"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of invalid nested object values", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"bar":      "string",
			"baz": map[interface{}]interface{}{
				"ban":  "int32",
				"bant": map[interface{}]interface{}{},
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorInvalidValue("bant", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})
}
