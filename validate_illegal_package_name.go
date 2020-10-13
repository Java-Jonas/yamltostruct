package yamltostruct

import (
	"fmt"
	"go/parser"
	"go/token"
)

func isValidPackageName(packageName string) bool {
	sourceCodeMock := `package ` + packageName

	_, err := parser.ParseFile(token.NewFileSet(), "", sourceCodeMock, 0)
	return err == nil
}

func validateIllegalPackageName(yamlData map[interface{}]interface{}) (errs []error) {

	for key, value := range yamlData {
		keyName := fmt.Sprintf("%v", key)

		if keyName == packageNameKey {
			packageName := fmt.Sprintf("%v", value)
			if !isValidPackageName(packageName) {
				errs = append(errs, newValidationErrorIllegalPackageName(packageName))
				return
			}
		}
	}

	return
}
