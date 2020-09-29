package yamltostruct

import (
	"fmt"
)

// returns errors if invalid values are used in the YAML file
// the declarations may not contain: Objects in Objects, Lists, "" and nil
func validateIllegalValue(yamlData map[interface{}]interface{}) (errs []error) {

	for key, value := range yamlData {
		keyName := fmt.Sprintf("%v", key)

		if isString(value) {
			if isEmptyString(value) {
				errs = append(errs, newValidationErrorIllegalValue(keyName, "root"))
			}
			continue
		}

		if isSlice(value) || isNil(value) {
			errs = append(errs, newValidationErrorIllegalValue(keyName, "root"))
			continue
		}

		if isMap(value) {
			mapValue, ok := value.(map[interface{}]interface{})
			if !ok {
				errs = append(errs, newUnexpectedError())
				continue
			}
			objectValidationErrs := validateIllegalValueObject(mapValue, keyName)
			errs = append(errs, objectValidationErrs...)
			continue
		}

		errs = append(errs, newValidationErrorIllegalValue(keyName, "root"))
	}

	return
}

func validateIllegalValueObject(yamlObjectData map[interface{}]interface{}, objectName string) (errs []error) {
	for key, value := range yamlObjectData {
		keyName := fmt.Sprintf("%v", key)

		if isString(value) {
			if isEmptyString(value) {
				errs = append(errs, newValidationErrorIllegalValue(keyName, objectName))
			}
			continue
		}

		if isSlice(value) || isMap(value) || isNil(value) {
			errs = append(errs, newValidationErrorIllegalValue(keyName, objectName))
			continue
		}

		errs = append(errs, newValidationErrorIllegalValue(keyName, objectName))
	}

	return
}
