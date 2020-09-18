package yamltostruct

import (
	"fmt"
	"reflect"
)

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
			objectValidationErrs := validateObject(mapValue, keyName)
			errs = append(errs, objectValidationErrs...)
			continue
		}

		errs = append(errs, newValidationErrorInvalidValue(keyName, "root"))
	}

	return
}

func validateObject(yamlObjectData map[interface{}]interface{}, parentObjectName string) (errs []error) {
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

func validateYamlData(yamlData map[interface{}]interface{}) (errs []error) {
	valueErrors := validateValues(yamlData)
	errs = append(errs, valueErrors...)

	return
}
