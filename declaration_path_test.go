package yamltostruct

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathBuilder(t *testing.T) {
	t.Run("should build segments with expected fieldLevels", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo": map[interface{}]interface{}{
				"bar": "string",
			},
		}

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "", data, fieldLevelZero)

		assert.Equal(t, 1, len(pb.paths))
		assert.Equal(t, []string{"foo.bar", "string"}, pb.paths[0].joinedNames())
		assert.Contains(t, pb.paths[0].segments, declaration{"foo", valueKindObject, firstFieldLevel})
		assert.Contains(t, pb.paths[0].segments, declaration{"bar", valueKindString, secondFieldLevel})
		assert.Contains(t, pb.paths[0].segments, declaration{"string", valueKindString, fieldLevelZero})
	})

	t.Run("should build one path of flat types and exit on loop", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "bar",
			"ban":      "foo",
			"bar":      "ban",
		}

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "ban", data["ban"], firstFieldLevel)

		assert.Equal(t, 1, len(pb.paths))
		assert.Equal(t, []string{"ban", "foo", "bar", "ban"}, pb.paths[0].joinedNames())
		assert.Equal(t, true, pb.paths[0].containsRecursiveness)
	})

	t.Run("should build one path of flat types that exits on a basic type", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "bar",
			"bar":      "string",
		}

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "foo", data["foo"], firstFieldLevel)

		assert.Equal(t, 1, len(pb.paths))
		assert.Equal(t, []string{"foo", "bar", "string"}, pb.paths[0].joinedNames())
		assert.Equal(t, false, pb.paths[0].containsRecursiveness)
	})

	t.Run("should build one path of a flat type directly", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar":      "string",
		}

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "bar", data["bar"], firstFieldLevel)

		assert.Equal(t, 1, len(pb.paths))
		assert.Equal(t, []string{"bar", "string"}, pb.paths[0].joinedNames())
		assert.Equal(t, false, pb.paths[0].containsRecursiveness)
	})

	t.Run("should build one path of nested types", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
			},
			"baz": "bar",
		}

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "baz", data["baz"], firstFieldLevel)

		assert.Equal(t, 1, len(pb.paths))
		assert.Equal(t, []string{"baz", "bar.foo", "baz"}, pb.paths[0].joinedNames())
		assert.Equal(t, true, pb.paths[0].containsRecursiveness)
	})

	t.Run("should build two paths of nested types", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
				"bam": "string",
			},
			"baz": "bar",
		}

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "baz", data["baz"], firstFieldLevel)

		assert.Equal(t, 2, len(pb.paths))
		joinedNamess := [][]string{
			pb.paths[0].joinedNames(),
			pb.paths[1].joinedNames(),
		}

		assert.Contains(t, joinedNamess, []string{"baz", "bar.foo", "baz"})
		assert.Contains(t, joinedNamess, []string{"baz", "bar.bam", "string"})
	})

	t.Run("should build a path with a itself-referring struct", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "bar",
			},
		}

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "bar", data["bar"], firstFieldLevel)

		assert.Equal(t, 1, len(pb.paths))
		assert.Equal(t, []string{"bar.foo", "bar"}, pb.paths[0].joinedNames())
		assert.Equal(t, true, pb.paths[0].containsRecursiveness)
	})

	t.Run("should build multiple paths of nested types", func(t *testing.T) {
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

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "baz", data["baz"], firstFieldLevel)

		assert.Equal(t, 5, len(pb.paths))
		joinedNamess := [][]string{
			pb.paths[0].joinedNames(),
			pb.paths[1].joinedNames(),
			pb.paths[2].joinedNames(),
			pb.paths[3].joinedNames(),
			pb.paths[4].joinedNames(),
		}

		assert.Contains(t, joinedNamess, []string{"baz", "bar.foo", "baz"})
		assert.Contains(t, joinedNamess, []string{"baz", "bar.bam", "string"})
		assert.Contains(t, joinedNamess, []string{"baz", "bar.bal", "bar"})
		assert.Contains(t, joinedNamess, []string{"baz", "bar.fof", "bas.ban", "string"})
		assert.Contains(t, joinedNamess, []string{"baz", "bar.fof", "bas.bunt", "bant", "int"})
	})

	t.Run("should build paths from yamlData root", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
				"bam": "string",
			},
			"baz": "bar",
		}

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "", data, fieldLevelZero)

		assert.Equal(t, 4, len(pb.paths))
		joinedNamess := [][]string{
			pb.paths[0].joinedNames(),
			pb.paths[1].joinedNames(),
			pb.paths[2].joinedNames(),
			pb.paths[3].joinedNames(),
		}

		assert.Contains(t, joinedNamess, []string{"baz", "bar.foo", "baz"})
		assert.Contains(t, joinedNamess, []string{"baz", "bar.bam", "string"})
		assert.Contains(t, joinedNamess, []string{"bar.bam", "string"})
		assert.Contains(t, joinedNamess, []string{"bar.foo", "baz", "bar"})
	})

	t.Run("should build recursive path from self referencing type", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"foo":      "foo",
		}

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "", data, fieldLevelZero)

		assert.Equal(t, 1, len(pb.paths))
		joinedNamess := [][]string{
			pb.paths[0].joinedNames(),
		}

		assert.Contains(t, joinedNamess, []string{"foo", "foo"})
		assert.Equal(t, pb.paths[0].containsRecursiveness, true)
	})

	t.Run("should build 2 recursive paths on self referencing objects", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"bar": map[interface{}]interface{}{
				"foo": "baz",
			},
			"baz": map[interface{}]interface{}{
				"ban": "bar",
			},
		}

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "", data, fieldLevelZero)

		assert.Equal(t, 2, len(pb.paths))
		joinedNamess := [][]string{
			pb.paths[0].joinedNames(),
			pb.paths[1].joinedNames(),
		}

		assert.Contains(t, joinedNamess, []string{"bar.foo", "baz.ban", "bar"})
		assert.Contains(t, joinedNamess, []string{"baz.ban", "bar.foo", "baz"})
		assert.Equal(t, pb.paths[0].containsRecursiveness, true)
		assert.Equal(t, pb.paths[1].containsRecursiveness, true)
	})

	t.Run("should stop paths on reference type used", func(t *testing.T) {
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

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "ban", data["ban"], firstFieldLevel)

		assert.Equal(t, 3, len(pb.paths))
		joinedNamess := [][]string{
			pb.paths[0].joinedNames(),
			pb.paths[1].joinedNames(),
			pb.paths[2].joinedNames(),
		}

		assert.Contains(t, joinedNamess, []string{"ban", "bar.foo", "*baz"})
		assert.Contains(t, joinedNamess, []string{"ban", "bar.fan", "[]baz"})
		assert.Contains(t, joinedNamess, []string{"ban", "bar.faz", "map[int]baz"})
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

		pb := pathBuilder{
			paths:    []declarationPath{},
			yamlData: data,
		}

		pb.build(declarationPath{}, "", data, fieldLevelZero)

		assert.Equal(t, 2, len(pb.paths))
		joinedNamess := [][]string{
			pb.paths[0].joinedNames(),
			pb.paths[1].joinedNames(),
		}

		assert.Contains(t, joinedNamess, []string{"bar.foo", "baz.ban", "bar"})
		assert.Contains(t, joinedNamess, []string{"baz.ban", "bar.foo", "baz"})
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

// 		pb := pathBuilder{
// 			paths: []declarationPath{},
// 			yamlData: data,
// 		}

// 		pb.build(declarationPath{}, "", data, fieldLevelZero)

// 		assert.Equal(t, 4, len(pb.paths))
// 		joinedNamess := [][]string{
// 			pb.paths[0].joinedNames(),
// 			pb.paths[1].joinedNames(),
// 			pb.paths[2].joinedNames(),
// 			pb.paths[3].joinedNames(),
// 		}

// 		assert.Contains(t, joinedNamess, []string{"foo", "int"})
// 		assert.Equal(t, pb.paths[0].containsRecursiveness, true)
// 	})
// }
