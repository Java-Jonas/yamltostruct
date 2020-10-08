package yamltostruct

import (
	"fmt"
)

type pathClosureKind int

const (
	// path ends due to it being recusrive
	pathClosureKindRecursiveness pathClosureKind = iota
	// path ends due to encountering a reference value
	pathClosureKindReference
)

type typeKind int

const (
	// eg. string, int
	typeKindValue typeKind = iota
	// eg. *string, map[int]string, []int
	typeKindReference
)

type yamlValueKind int

const (
	valueKindString yamlValueKind = iota
	valueKindObject
)

type fieldLevelKind int

const (
	// current level is the object itself
	fieldLevelZero fieldLevelKind = iota
	// current level is the root level of the object
	firstFieldLevel
	secondFieldLevel
)

// name it to be more like golang.ast package naming
type declaration struct {
	keyName       string
	yamlValueKind yamlValueKind
	fieldLevel    fieldLevelKind
	typeKind      typeKind
}

type declarationPath struct {
	declarations          []declaration
	closureKind           pathClosureKind
	containsRecursiveness bool
}

func (path *declarationPath) addDeclaration(
	keyName string,
	yamlValueKind yamlValueKind,
	fieldLevel fieldLevelKind,
	value interface{},
) {
	// TODO: maybe this shouldn't be done here
	for _, declaration := range path.declarations {
		if declaration.keyName == keyName {
			path.containsRecursiveness = true
		}
	}
	// TODO: identify typekind
	path.declarations = append(path.declarations, declaration{keyName, yamlValueKind, fieldLevel, typeKindValue})
}

// we list the declarations' typeNames with some additional logic
// to be more explicit (paths to nested field names will
// be concatenated eg. "foo.bar")
func (path declarationPath) joinedNames() []string {
	var joinedNames []string

	var wasStructField bool
	var parentStructName string

	for _, declaration := range path.declarations {
		if declaration.yamlValueKind == valueKindObject {
			wasStructField = true
			parentStructName = declaration.keyName
			continue
		}

		if wasStructField && declaration.fieldLevel == secondFieldLevel {
			joinedNames = append(joinedNames, parentStructName+"."+declaration.keyName)
		} else {
			joinedNames = append(joinedNames, declaration.keyName)
		}
		wasStructField = false
		parentStructName = ""

	}

	if wasStructField {
		joinedNames = append(joinedNames, parentStructName)
	}

	return joinedNames
}

func (path declarationPath) copySelf() declarationPath {
	declarationsCopy := make([]declaration, len(path.declarations))
	copy(declarationsCopy, path.declarations)
	pathCopy := path
	pathCopy.declarations = declarationsCopy
	return pathCopy
}

type pathBuilder struct {
	yamlData map[interface{}]interface{}
	paths    []declarationPath
}

func newPathBuilder(yamlData map[interface{}]interface{}) *pathBuilder {
	return &pathBuilder{yamlData: yamlData}
}

func (pb *pathBuilder) addPath(path declarationPath) {
	pb.paths = append(pb.paths, path)
}

// a recursive function to travel through the yamlData, creating
// a different path for each path
func (pb *pathBuilder) build(path declarationPath, keyName string, value interface{}, fieldLevel fieldLevelKind) {

	if fieldLevel == firstFieldLevel && keyName == "_package" {
		return
	}

	if isString(value) {
		path.addDeclaration(keyName, valueKindString, fieldLevel, value)
		if path.containsRecursiveness {
			// detected recursiveness implies this is the end of the path
			pb.addPath(path)
			return
		}
		valueLiteral := fmt.Sprintf("%v", value)
		// if a key cannot be found at root level we assume it's a basic type eg. string, int ..
		nextValue, isNotBasicType := pb.yamlData[valueLiteral]

		/*
			if isReferencingValue(){
				path.addDeclaration(valueLiteral, valueKindString, fieldLevelZero, value)
				// a reference type implies this is the end of the path
				pb.addPath(path)
				return
			}

		*/

		if !isNotBasicType {
			// TODO: add isBasicType func as *string would not be recognized as basic type
			path.addDeclaration(valueLiteral, valueKindString, fieldLevelZero, value)
			// a basic type implies this is the end of the path
			pb.addPath(path)
			return
		}
		// we get here only when the value is a reference to a user defined type
		// TODO: nextValue might also be *foo ... so ill have to find a different way (extract the actual type maybe)
		pb.build(path, valueLiteral, nextValue, firstFieldLevel)
	}

	if isMap(value) {
		path.addDeclaration(keyName, valueKindObject, fieldLevel, value)
		if path.containsRecursiveness {
			// detected recursiveness implies this is the end of the path
			pb.addPath(path)
			return
		}
		mapValue := value.(map[interface{}]interface{})
		for _key, _value := range mapValue {
			// the path is copied; this is basically a fork
			pathCopy := path.copySelf()
			_keyName := fmt.Sprintf("%v", _key)
			// we go a level deeper (fieldLevel+1) and handle each key/value
			// pair in the next pb.build() execution
			pb.build(pathCopy, _keyName, _value, fieldLevel+1)
		}
	}
}

func evalTypeKind(typeDefinitionString string) typeKind {
	return typeKindReference
}
