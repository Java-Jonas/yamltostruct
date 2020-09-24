package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateYamlRecursiveTypeUsage(t *testing.T) {
	t.Run("should not fail on usage of non-recursive types", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
			},
			"baz": map[interface{}]interface{}{
				"ban": "*bar",
			},
			"bal": map[interface{}]interface{}{
				"bam": "baz",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail when type is used in own declaration", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar":      "bar",
			"baz": map[interface{}]interface{}{
				"ban": "baz",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorRecursiveTypeUsage([]string{"bar"}),
			newValidationErrorRecursiveTypeUsage([]string{"baz.ban"}),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of recursive types (1/3)", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
			},
			"baz": map[interface{}]interface{}{
				"ban": "bar",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorRecursiveTypeUsage([]string{"bar.foo", "baz.ban"}),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of recursive types (2/3)", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "bam",
			},
			"baz": map[interface{}]interface{}{
				"ban": "bar",
			},
			"bam": map[interface{}]interface{}{
				"baf": "baz",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorRecursiveTypeUsage([]string{"bam.baf", "baz.ban", "bar.foo"}),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})
}

func TestDeclarationTreeGrow(t *testing.T) {
	t.Run("should build one branch of flat types and exit on loop", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "bar",
			"ban":      "foo",
			"bar":      "ban",
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "ban", data["ban"])

		assert.Equal(t, 1, len(dt.branches))
		assert.Equal(t, declarationBranch{"ban", "foo", "bar", "ban"}, dt.branches[0])
	})

	t.Run("should build one branch of flat types that exits on a basic type", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "bar",
			"bar":      "string",
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "foo", data["foo"])

		assert.Equal(t, 1, len(dt.branches))
		assert.Equal(t, declarationBranch{"foo", "bar", "string"}, dt.branches[0])
	})

	t.Run("should build one branch of a flat type directly", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar":      "string",
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "bar", data["bar"])

		assert.Equal(t, 1, len(dt.branches))
		assert.Equal(t, declarationBranch{"bar", "string"}, dt.branches[0])
	})

	t.Run("should build one branch of a flat type directly", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
			},
			"baz": "bar",
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "baz", data["baz"])

		assert.Equal(t, 1, len(dt.branches))
		assert.Equal(t, declarationBranch{"baz", "bar.", "foo", "baz"}, dt.branches[0])
	})
}
