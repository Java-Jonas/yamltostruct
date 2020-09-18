package yamltostruct

import (
	"fmt"
	"reflect"
	"regexp"
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

func validateValues(yamlData map[interface{}]interface{}) (errs []error) {

	for key, value := range yamlData {
		keyName := fmt.Sprintf("%v", key)

		if isString(value) {
			continue
		}

		if isSlice(value) {
			errs = append(errs, newValidationErrorInvalidValue(keyName, "root"))
			continue
		}

		if isMap(value) {
			mapValue, ok := value.(map[interface{}]interface{})
			if !ok {
				errs = append(errs, newUnexpectedError())
				continue
			}
			objectValidationErrs := validateObjectValues(mapValue, keyName)
			errs = append(errs, objectValidationErrs...)
			continue
		}

		errs = append(errs, newValidationErrorInvalidValue(keyName, "root"))
	}

	return
}

func validateObjectValues(yamlObjectData map[interface{}]interface{}, parentObjectName string) (errs []error) {
	for key, value := range yamlObjectData {
		keyName := fmt.Sprintf("%v", key)

		if isString(value) {
			continue
		}

		if isSlice(value) || isMap(value) {
			errs = append(errs, newValidationErrorInvalidValue(keyName, parentObjectName))
			continue
		}

		errs = append(errs, newValidationErrorInvalidValue(keyName, parentObjectName))
	}

	return
}

// extracts all types which are defined in a type definition
// map[string]int => []string{"string", "int"}
// (?!map\b)\b\w+/g
func extractTypes(typeDefinitionString string) (extractedTypes []string) {
	re := regexp.MustCompile(`(?!map\b)\b\w+/g`)
	// TODO
	re.FindString(typeDefinitionString)
	return
}

func findUndefinedTypesIn(usedTypes, definedTypes []string) (undefinedTypes []string) {
	return
}

func validateTypeDeclarationCompleteness(yamlData map[interface{}]interface{}) (errs []error) {

	var definedTypes []string

	for key := range yamlData {
		keyName := fmt.Sprintf("%v", key)
		if keyName == "_package" {
			continue
		}
		definedTypes = append(definedTypes, keyName)
	}

	for key, value := range yamlData {
		keyName := fmt.Sprintf("%v", key)

		if keyName == "_package" {
			continue
		}

		if isString(value) {
			valueString := fmt.Sprintf("%v", value)
			extractedTypes := extractTypes(valueString)
			undefinedTypes := findUndefinedTypesIn(extractedTypes, definedTypes)
			for _, undefiundefinedType := range undefinedTypes {
				errs = append(errs, newValidationErrorTypeNotFound(undefiundefinedType, "root"))
			}
		}

		if isMap(value) {
			mapValue, ok := value.(map[interface{}]interface{})
			if !ok {
				errs = append(errs, newUnexpectedError())
				continue
			}
			objectValidationErrs := validateObjectTypesDeclarationCompleteness(mapValue, keyName, definedTypes)
			errs = append(errs, objectValidationErrs...)
		}
	}

	return
}

func validateObjectTypesDeclarationCompleteness(
	yamlObjectData map[interface{}]interface{},
	parentObjectName string,
	definedTypes []string,
) (errs []error) {

	for _, value := range yamlObjectData {
		valueString := fmt.Sprintf("%v", value)
		extractedTypes := extractTypes(valueString)
		undefinedTypes := findUndefinedTypesIn(extractedTypes, definedTypes)
		for _, undefiundefinedType := range undefinedTypes {
			errs = append(errs, newValidationErrorTypeNotFound(undefiundefinedType, parentObjectName))
		}
	}

	return
}

func validatePackageDeclarationExistence(yamlData map[interface{}]interface{}) (errs []error) {
	var packageNameFound bool

	for key, value := range yamlData {
		keyName := fmt.Sprintf("%v", key)

		if keyName == "_package" && isString(value) {
			packageNameFound = true
		}
	}

	if !packageNameFound {
		errs = append(errs, newValidationErrorMissingPackageName())
	}

	return
}

func validateYamlData(yamlData map[interface{}]interface{}) (errs []error) {
	valueErrors := validateValues(yamlData)
	errs = append(errs, valueErrors...)

	missingPackageDeclarationErrs := validatePackageDeclarationExistence(yamlData)
	errs = append(errs, missingPackageDeclarationErrs...)

	return
}
