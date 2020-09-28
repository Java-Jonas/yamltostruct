package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateYamlMissingPackageName(t *testing.T) {
	t.Run("should not fail when the package name is specified", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"baz": map[interface{}]interface{}{
				"ban": "int",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail when the package name is not specified at root level", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"foo": "int",
			"baz": map[interface{}]interface{}{
				"ban":      "int",
				"_package": "foo",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{newValidationErrorMissingPackageName()}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail when the package name is not specified", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"foo": "int",
			"baz": map[interface{}]interface{}{
				"ban": "int",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{newValidationErrorMissingPackageName()}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})
}
