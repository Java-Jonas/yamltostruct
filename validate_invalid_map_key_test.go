package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
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
			"buf":      "map[map[int]bool]string",
			"baz": map[interface{}]interface{}{
				"ban": "map[[]foo]int",
			},
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{
			newValidationErrorInvalidMapKey("*foo", "map[*foo]int"),
			newValidationErrorInvalidMapKey("map[int]bool", "map[map[int]bool]string"),
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

	t.Run("should fail on usage of reference type as map key in nested map", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "[]string",
			"bar":      "map[int]map[foo]int",
			"baz": map[interface{}]interface{}{
				"bal": "map[bar]int",
			},
		}

		actualErrors := logicalValidation(data)
		expectedErrors := []error{
			newValidationErrorInvalidMapKey("foo", "map[int]map[foo]int"),
			newValidationErrorInvalidMapKey("bar", "map[bar]int"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})
}

// func TestExtractMapKeys(t *testing.T) {
// 	t.Run("should extract map keys from value strings", func(t *testing.T) {
// 		assert.Equal(t, extractMapKeys("map[int]string"), []string{"int"})
// 		assert.Equal(t, extractMapKeys("map[*int]string"), []string{"*int"})
// 		assert.Equal(t, extractMapKeys("map[[]int]string"), []string{"[]int"})
// 		assert.Equal(t, extractMapKeys("map[map[map[bool]int]string]float"), []string{"map[map[bool]int]string", "map[bool]int", "bool"})
// 		// assert.Equal(t, extractMapKeys("map[int]map[float]map[string]bool"), []string{"int", "float", "string"})
// 	})
// }

// func TestExtractRootLevelMapDeclarations(t *testing.T) {
// 	t.Run("should extract root level map declarations", func(t *testing.T) {
// 		assert.Equal(t, extractMapKeys("map[int]string"), []string{"map[int]string"})
// 	})
// }

func TestExtracpMapDeclExpression(t *testing.T) {
	t.Run("should extract map decl", func(t *testing.T) {
		mockSrc := `
	package main
	type mockType map[int]string
	`

		file, _ := parser.ParseFile(token.NewFileSet(), "", mockSrc, 0)
		assert.Equal(t, extracpMapDeclExpression(file).(*ast.MapType).Key.(*ast.Ident).Name, "int")
	})
}
