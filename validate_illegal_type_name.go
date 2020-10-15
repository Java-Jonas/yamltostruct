package yamltostruct

import (
	"fmt"
	"go/parser"
	"go/token"
)

// returns errors if type names contain illegal characters that do not adhere to golangs syntax restrictions
func validateIllegalTypeName(yamlData map[interface{}]interface{}) (errs []error) {
	for key, value := range yamlData {
		keyName := fmt.Sprintf("%v", key)

		if isIllegalTypeName(keyName) {
			errs = append(errs, newValidationErrorIllegalTypeName(keyName, "root"))
		}

		if isMap(value) {
			mapValue := value.(map[interface{}]interface{})
			objectValidationErrs := validateIllegalTypeNameObject(mapValue, keyName)
			errs = append(errs, objectValidationErrs...)
		}
	}

	return
}

func validateIllegalTypeNameObject(yamlObjectData map[interface{}]interface{}, objectName string) (errs []error) {
	for key := range yamlObjectData {
		keyName := fmt.Sprintf("%v", key)

		if isIllegalTypeName(keyName) {
			errs = append(errs, newValidationErrorIllegalTypeName(keyName, objectName))
		}
	}
	return
}

func isIllegalTypeName(typeName string) bool {
	sourceCodeMock := `
	package main
	type ` + typeName + ` string`

	_, err := parser.ParseFile(token.NewFileSet(), "", sourceCodeMock, 0)
	return err != nil
}
