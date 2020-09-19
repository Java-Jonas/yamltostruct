package yamltostruct

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findErrorIn(val error, slice []error) (int, bool) {
	for i, item := range slice {
		if item.Error() == val.Error() {
			return i, true
		}
	}
	return -1, false
}

func removeErrorFromSlice(slice []error, index int) []error {
	return append(slice[:index], slice[index+1:]...)
}

func matchErrors(actualErrors, expectedErrors []error) (leftoverErrors, redundantErrors []error) {
	// redefine redunantErrors so it never returns as nil (which happens when there are no redunant errors)
	// so comparing error slices becomes more conventient
	redundantErrors = make([]error, 0)
	leftoverErrors = make([]error, len(expectedErrors))
	copy(leftoverErrors, expectedErrors)

	for _, actualError := range actualErrors {
		leftoverErrorIndex, isFound := findErrorIn(actualError, leftoverErrors)
		if isFound {
			leftoverErrors = removeErrorFromSlice(leftoverErrors, leftoverErrorIndex)
		} else {
			redundantErrors = append(redundantErrors, actualError)
		}
	}

	return
}

func TestMatchErrors(t *testing.T) {
	t.Run("should return no errors when all errors matched", func(t *testing.T) {
		actualErrors := []error{errors.New("abc"), errors.New("def"), errors.New("ghi")}
		expectedErrors := []error{errors.New("abc"), errors.New("def"), errors.New("ghi")}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, 0, len(missingErrors))
		assert.Equal(t, 0, len(redundantErrors))
	})

	t.Run("should return missing errors when some are missing", func(t *testing.T) {
		actualErrors := []error{errors.New("abc")}
		expectedErrors := []error{errors.New("abc"), errors.New("def"), errors.New("ghi")}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		expectedMissingErrors := []error{errors.New("def"), errors.New("ghi")}
		assert.Equal(t, expectedMissingErrors, missingErrors)
		expectedRedundantErrors := []error{}
		assert.Equal(t, expectedRedundantErrors, redundantErrors)
	})

	t.Run("should return redundant errors when some are missing", func(t *testing.T) {
		actualErrors := []error{errors.New("abc"), errors.New("def"), errors.New("ghi")}
		expectedErrors := []error{errors.New("abc")}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		expectedMissingErrors := []error{}
		assert.Equal(t, expectedMissingErrors, missingErrors)
		expectedRedundantErrors := []error{errors.New("def"), errors.New("ghi")}
		assert.Equal(t, expectedRedundantErrors, redundantErrors)
	})

	t.Run("output should not be index order", func(t *testing.T) {
		actualErrors := []error{errors.New("abc"), errors.New("def"), errors.New("ghi")}
		expectedErrors := []error{errors.New("def"), errors.New("ghi"), errors.New("abc")}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		expectedMissingErrors := []error{}
		assert.Equal(t, expectedMissingErrors, missingErrors)
		expectedRedundantErrors := []error{}
		assert.Equal(t, expectedRedundantErrors, redundantErrors)
	})
}

func TestValidateYamlDataIllegalTypeName(t *testing.T) {
	t.Run("should not fail on valid key inputs", func(t *testing.T) {
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

	t.Run("should fail on spaces in key literal", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"fo o":     "int",
			"baz": map[interface{}]interface{}{
				"oof":  "int",
				"ba n": "int",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorIllegalTypeName("fo o", "root"),
			newValidationErrorIllegalTypeName("ba n", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should not fail on usage of allowed type names", func(t *testing.T) {
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

	t.Run("should fail on usage of keywords", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"break":    "int",
			"bar":      "string",
			"baz": map[interface{}]interface{}{
				"const": "int32",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorIllegalTypeName("break", "root"),
			newValidationErrorIllegalTypeName("const", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage special characters", func(t *testing.T) {
		data := map[interface{}]interface{}{
			"_package": "packageName",
			"*":        "int",
			"<":        "string",
			"fo$o":     "int",
			"baz": map[interface{}]interface{}{
				">-":    "int32",
				"bent{": "int32",
			},
		}

		actualErrors := validateYamlData(data)
		expectedErrors := []error{
			newValidationErrorIllegalTypeName("*", "root"),
			newValidationErrorIllegalTypeName("<", "root"),
			newValidationErrorIllegalTypeName("fo$o", "root"),
			newValidationErrorIllegalTypeName(">-", "baz"),
			newValidationErrorIllegalTypeName("bent{", "baz"),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})
}

func TestIsIllegalTypeName(t *testing.T) {
	t.Run("should return false if the type names are valid", func(t *testing.T) {
		assert.Equal(t, false, isIllegalTypeName("foo_"), isIllegalTypeName("b_ar"), isIllegalTypeName("BA2Z"))
	})
	t.Run("should return true if the type names are illegal", func(t *testing.T) {
		assert.Equal(t, true, isIllegalTypeName("fo o"), isIllegalTypeName("b*ar"), isIllegalTypeName("B+2Z"))
	})
}

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

		actualErrors := validateYamlData(data)
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

		actualErrors := validateYamlData(data)
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

		actualErrors := validateYamlData(data)
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

		actualErrors := validateYamlData(data)
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

		actualErrors := validateYamlData(data)
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

		actualErrors := validateYamlData(data)
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

		actualErrors := validateYamlData(data)
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

		actualErrors := validateYamlData(data)
		expectedErrors := []error{}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

}

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
			newValidationErrorRecursiveTypeUsage([]string{"ban"}),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of recursive types (1/2)", func(t *testing.T) {
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
			newValidationErrorRecursiveTypeUsage([]string{"foo", "ban"}),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})

	t.Run("should fail on usage of recursive types (2/2)", func(t *testing.T) {
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
			newValidationErrorRecursiveTypeUsage([]string{"baf", "ban", "foo"}),
		}

		missingErrors, redundantErrors := matchErrors(actualErrors, expectedErrors)

		assert.Equal(t, []error{}, missingErrors)
		assert.Equal(t, []error{}, redundantErrors)
	})
}

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

func TestExtractTypes(t *testing.T) {
	t.Run("should extract a basic type", func(t *testing.T) {
		input := "string"

		actualOutput := extractTypes(input)
		expectedOutput := []string{"string"}

		assert.Equal(t, expectedOutput, actualOutput)
	})
	t.Run("should extract a slice type", func(t *testing.T) {
		input := "[]int"

		actualOutput := extractTypes(input)
		expectedOutput := []string{"int"}

		assert.Equal(t, expectedOutput, actualOutput)
	})
	t.Run("should extract both types form a map declaration", func(t *testing.T) {
		input := "map[string]int16"

		actualOutput := extractTypes(input)
		expectedOutput := []string{"string", "int16"}

		assert.Equal(t, expectedOutput, actualOutput)
	})
	t.Run("should extract all types from a complicated declaration", func(t *testing.T) {
		input := "map[string]map[uint][][]bool"

		actualOutput := extractTypes(input)
		expectedOutput := []string{"string", "uint", "bool"}

		assert.Equal(t, expectedOutput, actualOutput)
	})
}

func TestFindUndefinedTypesIn(t *testing.T) {
	t.Run("should find all undefined types", func(t *testing.T) {
		definedTypesInput := []string{"foo", "bar"}
		usedTypesInput := []string{"foo", "bar", "baz", "string", "uint16", "bool", "bam"}

		actualOutput := findUndefinedTypesIn(usedTypesInput, definedTypesInput)
		expectedOutput := []string{"baz", "bam"}

		assert.Equal(t, expectedOutput, actualOutput)
	})
}
