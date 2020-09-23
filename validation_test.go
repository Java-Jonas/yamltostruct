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

func TestIsIllegalTypeName(t *testing.T) {
	t.Run("should return false if the type names are valid", func(t *testing.T) {
		assert.Equal(t, false, isIllegalTypeName("foo_"), isIllegalTypeName("b_ar"), isIllegalTypeName("BA2Z"))
	})
	t.Run("should return true if the type names are illegal", func(t *testing.T) {
		assert.Equal(t, true, isIllegalTypeName("fo o"), isIllegalTypeName("b*ar"), isIllegalTypeName("B+2Z"))
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
