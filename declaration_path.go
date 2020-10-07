package yamltostruct

import (
	"fmt"
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

type declaration struct {
	keyName       string
	yamlValueKind yamlValueKind
	// TODO typeKind value/reference
	fieldLevel fieldLevelKind
}

type declarationPath struct {
	segments []declaration
	// TODO: pathClosureKind (recursiveness, reference type ..)
	containsRecursiveness bool
}

func (path *declarationPath) addSegment(keyName string, yamlValueKind yamlValueKind, fieldLevel fieldLevelKind) {
	// TODO: maybe this shouldn't be done here
	for _, segment := range path.segments {
		if segment.keyName == keyName {
			path.containsRecursiveness = true
		}
	}
	path.segments = append(path.segments, declaration{keyName, yamlValueKind, fieldLevel})
}

// we list the segments' typeNames with some additional logic
// to be more explicit (paths to nested field names will
// be concatenated eg. "foo.bar")
func (path declarationPath) joinedNames() []string {
	var joinedNames []string

	var wasStructField bool
	var parentStructName string

	for _, segment := range path.segments {
		if segment.yamlValueKind == valueKindObject {
			wasStructField = true
			parentStructName = segment.keyName
			continue
		}

		if wasStructField && segment.fieldLevel == secondFieldLevel {
			joinedNames = append(joinedNames, parentStructName+"."+segment.keyName)
		} else {
			joinedNames = append(joinedNames, segment.keyName)
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
	segmentsCopy := make([]declaration, len(path.segments))
	copy(segmentsCopy, path.segments)
	pathCopy := path
	pathCopy.segments = segmentsCopy
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
		path.addSegment(keyName, valueKindString, fieldLevel)
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
				path.addSegment(valueLiteral, valueKindString, fieldLevelZero)
				// a reference type implies this is the end of the path
				pb.addPath(path)
				return
			}

		*/

		if !isNotBasicType {
			path.addSegment(valueLiteral, valueKindString, fieldLevelZero)
			// a basic type implies this is the end of the path
			pb.addPath(path)
			return
		}
		// we get here only when the value is a reference to a user defined type
		pb.build(path, valueLiteral, nextValue, firstFieldLevel)
	}

	if isMap(value) {
		path.addSegment(keyName, valueKindObject, fieldLevel)
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
