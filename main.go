package yamltostruct

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

func readYamlFile(filePath string) ([]byte, error) {
	yamlDataBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		return yamlDataBytes, err
	}

	return yamlDataBytes, err
}

func ConvertToDataMap(yamlDataBytes []byte) (map[interface{}]interface{}, error) {
	yamlData := make(map[interface{}]interface{})
	err := yaml.Unmarshal(yamlDataBytes, &yamlData)

	if err != nil {
		return yamlData, err
	}

	return yamlData, err
}

func Unmarshal(yamlDataBytes []byte) ([]byte, error) {
	return []byte("test"), nil
}

func newValidationErrorTypeNotFound(missingTypeLiteral, parentItemName string) error {
	return errors.New(
		fmt.Sprintf(
			"ErrTypeNotFound: type with name \"%s\" in \"%s\" was not found",
			missingTypeLiteral,
			parentItemName,
		),
	)
}
func newValidationErrorMissingPackageName() error {
	return errors.New("ErrMissingPackageName: package name was not specified in the \"_package\" field at root level")
}
func newUnexpectedError() error {
	return errors.New("an unexpected error occured")
}
func newValidationErrorInvalidValue(keyName, parentItemName string) error {
	return errors.New(
		fmt.Sprintf(
			"ErrInvalidType: value assigned to key \"%s\" in \"%s\" is invalid",
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

type ASTBuilder struct{}

func NewASTBuilder() *ASTBuilder {
	return &ASTBuilder{}
}

func (a *ASTBuilder) build(yamlData map[interface{}]interface{}) *ast.File {
	src := `
	package main
	`
	f, _ := parser.ParseFile(token.NewFileSet(), "", src, 0)
	return f
}
