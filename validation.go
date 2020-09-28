package yamltostruct

import (
	"fmt"
	"reflect"
)

var golangBasicTypes = []string{"string", "bool", "int8", "uint8", "byte", "int16", "uint16", "int32", "rune", "uint32", "int64", "uint64", "int", "uint", "uintptr", "float32", "float64", "complex64", "complex128"}

func isString(unknown interface{}) bool {
	v := reflect.ValueOf(unknown)
	if v.Kind() == reflect.String {
		return true
	}
	return false
}

func isSlice(unknown interface{}) bool {
	v := reflect.ValueOf(unknown)
	if v.Kind() == reflect.Slice {
		return true
	}
	return false
}

func isMap(unknown interface{}) bool {
	v := reflect.ValueOf(unknown)
	if v.Kind() == reflect.Map {
		return true
	}
	return false
}

func isNil(unknown interface{}) bool {
	return unknown == nil
}

func isEmptyString(unknown interface{}) bool {
	if !isString(unknown) {
		return false
	}
	valueString := fmt.Sprintf("%v", unknown)
	return valueString == ""
}

func validateYamlData(yamlData map[interface{}]interface{}) (errs []error) {
	valueErrors := validateIllegalValue(yamlData)
	errs = append(errs, valueErrors...)

	illegalTypeNameErrors := validateIllegalTypeName(yamlData)
	errs = append(errs, illegalTypeNameErrors...)

	missingPackageDeclarationErrs := validateMissingPackageName(yamlData)
	errs = append(errs, missingPackageDeclarationErrs...)

	missingTypeDeclarationErrs := validateTypeNotFound(yamlData)
	errs = append(errs, missingTypeDeclarationErrs...)

	recursiveTypeUsageErrs := validateRecursiveTypeUsage(yamlData)
	errs = append(errs, recursiveTypeUsageErrs...)

	return
}
