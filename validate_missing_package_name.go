package yamltostruct

import (
	"fmt"
)

// returns errors if a package name was not declared with the key "_package" in the YAML file
func validateMissingPackageName(yamlData map[interface{}]interface{}) (errs []error) {
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
