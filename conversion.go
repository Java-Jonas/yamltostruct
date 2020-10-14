package yamltostruct

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func convertToAST(yamlData map[interface{}]interface{}) *ast.File {
	sw := &sourceWriter{}

	for key, value := range yamlData {
		keyName := fmt.Sprintf("%v", key)

		if isString(value) {
			valueString := fmt.Sprintf("%v", value)
			if keyName == packageNameKey {
				sw.addPackageName(valueString)
			} else {
				sw.addNamedType(keyName, valueString)
			}
			continue
		}

		if isMap(value) {
			mapValue := value.(map[interface{}]interface{})
			sw.startStructType(keyName)
			for _key, _value := range mapValue {
				_valueString := fmt.Sprintf("%v", _value)
				_keyName := fmt.Sprintf("%v", _key)
				sw.addStructField(_keyName, _valueString)
			}
			sw.closeStructType()
		}
	}

	return sw.parse()
}

type sourceWriter struct {
	sourceCode string
}

func (s *sourceWriter) parse() *ast.File {
	file, _ := parser.ParseFile(token.NewFileSet(), "", s.sourceCode, 0)
	return file
}

func (s *sourceWriter) addPackageName(packageName string) *sourceWriter {
	s.sourceCode = fmt.Sprintf("package %s\n%s", packageName, s.sourceCode)
	return s
}

func (s *sourceWriter) addNamedType(name, typeName string) *sourceWriter {
	s.sourceCode = fmt.Sprintf("%s\ntype %s %s", s.sourceCode, name, typeName)
	return s
}

func (s *sourceWriter) startStructType(name string) *sourceWriter {
	s.sourceCode = fmt.Sprintf("%s\ntype %s struct {", s.sourceCode, name)
	return s
}

func (s *sourceWriter) addStructField(name, typeName string) *sourceWriter {
	s.sourceCode = fmt.Sprintf("%s\n%s %s", s.sourceCode, name, typeName)
	return s
}

func (s *sourceWriter) closeStructType() *sourceWriter {
	s.sourceCode = fmt.Sprintf("%s\n}", s.sourceCode)
	return s
}
