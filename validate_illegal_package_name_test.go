package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateYamlIllegalPackageName(t *testing.T) {
	t.Run("should not fail when the package name is valid", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packagename",
		}

		actualErrors := syntacticalValidation(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail when the package name is not valid", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "package-name",
		}

		actualErrors := syntacticalValidation(data)
		expectedErrors := []error{
			newValidationErrorIllegalPackageName("package-name"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})
}
