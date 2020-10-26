package yamltostruct

import (
	"errors"
	"fmt"
	"strings"
)

func newValidationErrorTypeNotFound(missingTypeLiteral, parentItemName string) error {
	return errors.New(
		fmt.Sprintf(
			"ErrTypeNotFound: type with name \"%s\" in \"%s\" was not found",
			missingTypeLiteral,
			parentItemName,
		),
	)
}
func newValidationErrorIllegalValue(keyName, parentItemName string) error {
	return errors.New(
		fmt.Sprintf(
			"ErrIllegalValue: value assigned to key \"%s\" in \"%s\" is invalid",
			keyName,
			parentItemName,
		),
	)
}
func newValidationErrorInvalidValueString(valueString, keyName, parentItemName string) error {
	return errors.New(
		fmt.Sprintf(
			"ErrInvalidValueString: value \"%s\" assigned to \"%s\" in \"%s\" is invalid",
			valueString,
			keyName,
			parentItemName,
		),
	)
}
func newValidationErrorIllegalTypeName(keyName, parentItemName string) error {
	return errors.New(
		fmt.Sprintf(
			"ErrIllegalTypeName: illegal type name \"%s\" in \"%s\"",
			keyName,
			parentItemName,
		),
	)
}
func newValidationErrorRecursiveTypeUsage(keysResultingInRecursiveness []string) error {
	keys := strings.Join(keysResultingInRecursiveness, "->")
	return errors.New(
		fmt.Sprintf(
			"ErrRecursiveTypeUsage: illegal recursive type detected for \"%s\"",
			keys,
		),
	)
}
func newValidationErrorInvalidMapKey(mapKey, valueString string) error {
	return errors.New(
		fmt.Sprintf(
			"ErrInvalidMapKey: \"%s\" in \"%s\" is not a valid map key",
			mapKey,
			valueString,
		),
	)
}
