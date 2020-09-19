package yamltostruct

import (
	"fmt"
	"regexp"
)

var golangKeywords []string = []string{
	"break", "default", "func", "interface", "select", "case",
	"defer", "go", "map", "struct", "chan", "else", "goto",
	"package", "switch", "const", "fallthrough", "if", "range",
	"type", "continue", "for", "import", "return", "var",
}

// returns errors if type names contain illegal characters that do not adhere to golangs syntax restrictions
func validateTypeNames(yamlData map[interface{}]interface{}) (errs []error) {
	for key, value := range yamlData {
		keyName := fmt.Sprintf("%v", key)

		if isIllegalTypeName(keyName) {
			errs = append(errs, newValidationErrorIllegalTypeName(keyName, "root"))
		}

		if isMap(value) {
			mapValue, ok := value.(map[interface{}]interface{})
			if !ok {
				errs = append(errs, newUnexpectedError())
				continue
			}
			objectValidationErrs := validateObjectFieldNames(mapValue, keyName)
			errs = append(errs, objectValidationErrs...)
		}
	}

	return
}

func validateObjectFieldNames(yamlObjectData map[interface{}]interface{}, objectName string) (errs []error) {
	for key := range yamlObjectData {
		keyName := fmt.Sprintf("%v", key)

		if isIllegalTypeName(keyName) {
			errs = append(errs, newValidationErrorIllegalTypeName(keyName, objectName))
		}
	}
	return
}

func isKeyword(literal string) bool {
	for _, keyword := range golangKeywords {
		if literal == keyword {
			return true
		}
	}
	return false
}

func isIllegalTypeName(typeName string) bool {
	re := regexp.MustCompile(`[^A-Za-z0-9_]`)
	isIllegal := re.MatchString(typeName)
	isKeyword := isKeyword(typeName)
	return isIllegal || isKeyword
}

/* TODO: keywords
 */
