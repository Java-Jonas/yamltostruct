package yamltostruct

import (
	"fmt"
)

// returns errors if invalid values are used in the YAML file
// the declarations may not contain: Objects in Objects, Lists
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

func validateObjectValues(yamlObjectData map[interface{}]interface{}, objectName string) (errs []error) {
	for key, value := range yamlObjectData {
		keyName := fmt.Sprintf("%v", key)

		if isString(value) {
			continue
		}

		if isSlice(value) || isMap(value) {
			errs = append(errs, newValidationErrorInvalidValue(keyName, objectName))
			continue
		}

		errs = append(errs, newValidationErrorInvalidValue(keyName, objectName))
	}

	return
}
