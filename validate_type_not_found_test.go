package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateYamlDataTypeNotFound(t *testing.T) {
	t.Run("should not fail on usage of standard types", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"bar":      "string",
			"baf":      "[]string",
			"bal":      "map[string]int",
			"baz": map[interface{}]interface{}{
				"ban":  "int32",
				"bunt": "[]int",
				"bap":  "map[int16]string",
			},
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should not fail on usage of declared types", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"bar":      "string",
			"baf":      "[]foo",
			"bal":      "map[foo]bar",
			"bum":      "*int",
			"baz": map[interface{}]interface{}{
				"ban":  "int32",
				"bam":  "bar",
				"bunt": "[]baf",
				"bap":  "map[bar]foo",
				"bal":  "***bar",
				"lap":  "map[**bar]**foo",
				"slap": "**[]**baf",
			},
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of types declared in object", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"baz": map[interface{}]interface{}{
				"ban": "int32",
				"bar": "ban",
			},
			"boo": "ban",
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{
			newValidationErrorTypeNotFound("ban", "baz"),
			newValidationErrorTypeNotFound("ban", "root"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of unknown types", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"fof":      "schtring",
			"baz": map[interface{}]interface{}{
				"ban": "int32",
				"bam": "bar",
			},
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{
			newValidationErrorTypeNotFound("schtring", "root"),
			newValidationErrorTypeNotFound("bar", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of unknown types in slices", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"fof":      "[]schtring",
			"baz": map[interface{}]interface{}{
				"ban": "int32",
				"bam": "[]bar",
			},
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{
			newValidationErrorTypeNotFound("schtring", "root"),
			newValidationErrorTypeNotFound("bar", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of unknown types in maps", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "int",
			"fof":      "[int]schtring",
			"boo":      "[schtring]int",
			"baz": map[interface{}]interface{}{
				"ban": "int32",
				"bam": "[int]bar",
				"bal": "[bar]int",
			},
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{
			newValidationErrorTypeNotFound("schtring", "root"),
			newValidationErrorTypeNotFound("schtring", "root"),
			newValidationErrorTypeNotFound("bar", "baz"),
			newValidationErrorTypeNotFound("bar", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail with multiple errors of multiple undefined types are used in declaration", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "map[bar]map[ban]baz",
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{
			newValidationErrorTypeNotFound("bar", "root"),
			newValidationErrorTypeNotFound("ban", "root"),
			newValidationErrorTypeNotFound("baz", "root"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should not fail when type is used before declared", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"fof":      "foo",
			"foo":      "int",
			"baz": map[interface{}]interface{}{
				"ban": "int32",
				"bam": "bar",
			},
			"bar": "string",
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

}
