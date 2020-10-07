package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeclarationTreeGrow(t *testing.T) {
	t.Run("should build segments with expected fieldLevels", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo": map[interface{}]interface{}{
				"bar": "string",
			},
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "", data, fieldLevelZero)

		assert.Equal(t, 1, len(dt.branches))
		assert.Equal(t, []string{"foo.bar", "string"}, dt.branches[0].declarationPath())
		assert.Contains(t, dt.branches[0].segments, declarationSegment{"foo", valueKindObject, firstFieldLevel})
		assert.Contains(t, dt.branches[0].segments, declarationSegment{"bar", valueKindString, secondFieldLevel})
		assert.Contains(t, dt.branches[0].segments, declarationSegment{"string", valueKindString, fieldLevelZero})
	})

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

		dt.grow(declarationBranch{}, "ban", data["ban"], firstFieldLevel)

		assert.Equal(t, 1, len(dt.branches))
		assert.Equal(t, []string{"ban", "foo", "bar", "ban"}, dt.branches[0].declarationPath())
		assert.Equal(t, true, dt.branches[0].containsRecursiveness)
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

		dt.grow(declarationBranch{}, "foo", data["foo"], firstFieldLevel)

		assert.Equal(t, 1, len(dt.branches))
		assert.Equal(t, []string{"foo", "bar", "string"}, dt.branches[0].declarationPath())
		assert.Equal(t, false, dt.branches[0].containsRecursiveness)
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

		dt.grow(declarationBranch{}, "bar", data["bar"], firstFieldLevel)

		assert.Equal(t, 1, len(dt.branches))
		assert.Equal(t, []string{"bar", "string"}, dt.branches[0].declarationPath())
		assert.Equal(t, false, dt.branches[0].containsRecursiveness)
	})

	t.Run("should build one branch of nested types", func(t *testing.T) {
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

		dt.grow(declarationBranch{}, "baz", data["baz"], firstFieldLevel)

		assert.Equal(t, 1, len(dt.branches))
		assert.Equal(t, []string{"baz", "bar.foo", "baz"}, dt.branches[0].declarationPath())
		assert.Equal(t, true, dt.branches[0].containsRecursiveness)
	})

	t.Run("should build two branches of nested types", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
				"bam": "string",
			},
			"baz": "bar",
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "baz", data["baz"], firstFieldLevel)

		assert.Equal(t, 2, len(dt.branches))
		declarationPaths := [][]string{
			dt.branches[0].declarationPath(),
			dt.branches[1].declarationPath(),
		}

		assert.Contains(t, declarationPaths, []string{"baz", "bar.foo", "baz"})
		assert.Contains(t, declarationPaths, []string{"baz", "bar.bam", "string"})
	})

	t.Run("should build a branch with a itself-referring struct", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "bar",
			},
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "bar", data["bar"], firstFieldLevel)

		assert.Equal(t, 1, len(dt.branches))
		assert.Equal(t, []string{"bar.foo", "bar"}, dt.branches[0].declarationPath())
		assert.Equal(t, true, dt.branches[0].containsRecursiveness)
	})

	t.Run("should build multiple branches of nested types", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
				"bam": "string",
				"bal": "bar",
				"fof": "bas",
			},
			"bas": map[interface{}]interface{}{
				"ban":  "string",
				"bunt": "bant",
			},
			"baz":  "bar",
			"bant": "int",
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "baz", data["baz"], firstFieldLevel)

		assert.Equal(t, 5, len(dt.branches))
		declarationPaths := [][]string{
			dt.branches[0].declarationPath(),
			dt.branches[1].declarationPath(),
			dt.branches[2].declarationPath(),
			dt.branches[3].declarationPath(),
			dt.branches[4].declarationPath(),
		}

		assert.Contains(t, declarationPaths, []string{"baz", "bar.foo", "baz"})
		assert.Contains(t, declarationPaths, []string{"baz", "bar.bam", "string"})
		assert.Contains(t, declarationPaths, []string{"baz", "bar.bal", "bar"})
		assert.Contains(t, declarationPaths, []string{"baz", "bar.fof", "bas.ban", "string"})
		assert.Contains(t, declarationPaths, []string{"baz", "bar.fof", "bas.bunt", "bant", "int"})
	})

	t.Run("should build branches from yamlData root", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
				"bam": "string",
			},
			"baz": "bar",
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "", data, fieldLevelZero)

		assert.Equal(t, 4, len(dt.branches))
		declarationPaths := [][]string{
			dt.branches[0].declarationPath(),
			dt.branches[1].declarationPath(),
			dt.branches[2].declarationPath(),
			dt.branches[3].declarationPath(),
		}

		assert.Contains(t, declarationPaths, []string{"baz", "bar.foo", "baz"})
		assert.Contains(t, declarationPaths, []string{"baz", "bar.bam", "string"})
		assert.Contains(t, declarationPaths, []string{"bar.bam", "string"})
		assert.Contains(t, declarationPaths, []string{"bar.foo", "baz", "bar"})
	})

	t.Run("should build recursive branch from self referencing type", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "foo",
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "", data, fieldLevelZero)

		assert.Equal(t, 1, len(dt.branches))
		declarationPaths := [][]string{
			dt.branches[0].declarationPath(),
		}

		assert.Contains(t, declarationPaths, []string{"foo", "foo"})
		assert.Equal(t, dt.branches[0].containsRecursiveness, true)
	})

	t.Run("should build 2 recursive branches on self referencing objects", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
			},
			"baz": map[interface{}]interface{}{
				"ban": "bar",
			},
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "", data, fieldLevelZero)

		assert.Equal(t, 2, len(dt.branches))
		declarationPaths := [][]string{
			dt.branches[0].declarationPath(),
			dt.branches[1].declarationPath(),
		}

		assert.Contains(t, declarationPaths, []string{"bar.foo", "baz.ban", "bar"})
		assert.Contains(t, declarationPaths, []string{"baz.ban", "bar.foo", "baz"})
		assert.Equal(t, dt.branches[0].containsRecursiveness, true)
		assert.Equal(t, dt.branches[1].containsRecursiveness, true)
	})

	t.Run("should stop branches on reference type used", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "*baz",
				"fan": "[]baz",
				"faz": "map[int]baz",
			},
			"baz": "int",
			"ban": "bar",
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "ban", data["ban"], firstFieldLevel)

		assert.Equal(t, 3, len(dt.branches))
		declarationPaths := [][]string{
			dt.branches[0].declarationPath(),
			dt.branches[1].declarationPath(),
			dt.branches[2].declarationPath(),
		}

		assert.Contains(t, declarationPaths, []string{"ban", "bar.foo", "*baz"})
		assert.Contains(t, declarationPaths, []string{"ban", "bar.fan", "[]baz"})
		assert.Contains(t, declarationPaths, []string{"ban", "bar.faz", "map[int]baz"})
	})
}

func TestDeclarationPath(t *testing.T) {
	t.Run("should build declaration path", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
			},
			"baz": map[interface{}]interface{}{
				"ban": "bar",
			},
		}

		dt := declarationTree{
			branches: []declarationBranch{},
			yamlData: data,
		}

		dt.grow(declarationBranch{}, "", data, fieldLevelZero)

		assert.Equal(t, 2, len(dt.branches))
		declarationPaths := [][]string{
			dt.branches[0].declarationPath(),
			dt.branches[1].declarationPath(),
		}

		assert.Contains(t, declarationPaths, []string{"bar.foo", "baz.ban", "bar"})
		assert.Contains(t, declarationPaths, []string{"baz.ban", "bar.foo", "baz"})
	})
}

// func TestTypeKind(t *testing.T) {
// 	t.Run("typeKind is detected for reference/value types", func(t *testing.T) {
// 		data := map[interface{}]interface{}{
// 			"_package": "packageName",
// 			"foo":      "int",
// 			"bar":      "*int",
// 			"ban":      "[]int",
// 			"baf":      "map[int]string",
// 		}

// 		dt := declarationTree{
// 			branches: []declarationBranch{},
// 			yamlData: data,
// 		}

// 		dt.grow(declarationBranch{}, "", data, fieldLevelZero)

// 		assert.Equal(t, 4, len(dt.branches))
// 		declarationPaths := [][]string{
// 			dt.branches[0].declarationPath(),
// 			dt.branches[1].declarationPath(),
// 			dt.branches[2].declarationPath(),
// 			dt.branches[3].declarationPath(),
// 		}

// 		assert.Contains(t, declarationPaths, []string{"foo", "int"})
// 		assert.Equal(t, dt.branches[0].containsRecursiveness, true)
// 	})
// }
