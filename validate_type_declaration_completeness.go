package yamltostruct

import (
	"fmt"
	"regexp"
)

// returns errors if types are used which are not declared in the YAML file
// order of declaration is irrelevant
func validateObjectTypesDeclarationCompleteness(
	yamlObjectData map[interface{}]interface{},
	parentObjectName string,
	definedTypes []string,
) (errs []error) {

	for _, value := range yamlObjectData {
		valueString := fmt.Sprintf("%v", value)
		extractedTypes := extractTypes(valueString)
		undefinedTypes := findUndefinedTypesIn(extractedTypes, definedTypes)
		for _, undefinedType := range undefinedTypes {
			errs = append(errs, newValidationErrorTypeNotFound(undefinedType, parentObjectName))
		}
	}

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
			for _, undefinedType := range undefinedTypes {
				errs = append(errs, newValidationErrorTypeNotFound(undefinedType, "root"))
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

// extracts all types which are defined in a type definition
// map[string]int => []string{"string", "int"}
// (?!map\b)\b\w+/g
func extractTypes(typeDefinitionString string) (extractedTypes []string) {
	re := regexp.MustCompile(`[A-Za-z0-9_]*`)
	matches := re.FindAllString(typeDefinitionString, -1)
	for _, match := range matches {
		if match == "map" || match == "" {
			continue
		}
		extractedTypes = append(extractedTypes, match)
	}
	return
}

func findUndefinedTypesIn(usedTypes, definedTypes []string) (undefinedTypes []string) {
	allKnownTypes := append(definedTypes, golangBasicTypes...)
	for _, usedType := range usedTypes {
		var isDefined bool
		for _, knownType := range allKnownTypes {
			if knownType == usedType {
				isDefined = true
				break
			}
		}
		if !isDefined {
			undefinedTypes = append(undefinedTypes, usedType)
		}
	}
	return
}
