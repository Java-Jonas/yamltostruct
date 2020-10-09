package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateYamlDataInvalidMapKey(t *testing.T) {
	t.Run("should not fail on usage of valid map keys", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "string",
			"bal":      "map[foo]int",
			"baz": map[interface{}]interface{}{
				"bal": "map[foo]int",
			},
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of reference type directly as map key", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "string",
			"bar":      "map[*foo]int",
			"baz": map[interface{}]interface{}{
				"ban": "map[[]foo]int",
			},
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{
			newValidationErrorInvalidMapKey("*foo", "map[*foo]int"),
			newValidationErrorInvalidMapKey("[]foo", "map[[]foo]int"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of reference type as map key", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "[]string",
			"ban":      "*int",
			"bunt":     "map[int]string",
			"bar":      "map[foo]int",
			"baz": map[interface{}]interface{}{
				"bal": "map[ban]int",
				"buf": "map[bunt]int",
			},
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{
			newValidationErrorInvalidMapKey("foo", "map[foo]int"),
			newValidationErrorInvalidMapKey("ban", "map[ban]int"),
			newValidationErrorInvalidMapKey("bunt", "map[bunt]int"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})
}
