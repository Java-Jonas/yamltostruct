package yamltostruct

import (
	"fmt"
)

type valueKind int

const (
	valueKindString valueKind = iota
	valueKindObject
)

type fieldLevelKind int

const (
	fieldLevelZero fieldLevelKind = iota
	firstFieldLevel
	secondFieldLevel
)

type declarationSegment struct {
	typeName   string
	valueType  valueKind
	fieldLevel fieldLevelKind
}

type declarationBranch struct {
	segments              []declarationSegment
	containsRecursiveness bool
}

func (branch *declarationBranch) addSegment(typeName string, valueType valueKind, fieldLevel fieldLevelKind) {
	for _, segment := range branch.segments {
		if segment.typeName == typeName {
			branch.containsRecursiveness = true
		}
	}
	branch.segments = append(branch.segments, declarationSegment{typeName, valueType, fieldLevel})
}

func (branch declarationBranch) declarationPath() []string {
	var path []string

	var isStructField bool
	var parentStructName string

	for _, segment := range branch.segments {
		if segment.valueType == valueKindObject {
			isStructField = true
			parentStructName = segment.typeName
			continue
		}
		if isStructField {
			path = append(path, parentStructName+"."+segment.typeName)
			isStructField = false
			parentStructName = ""
		} else {
			path = append(path, segment.typeName)
		}
	}

	return path
}

func (branch declarationBranch) copySelf() declarationBranch {
	segmentsCopy := make([]declarationSegment, len(branch.segments))
	copy(segmentsCopy, branch.segments)
	branchCopy := branch
	branchCopy.segments = segmentsCopy
	return branchCopy
}

type declarationTree struct {
	yamlData map[interface{}]interface{}
	branches []declarationBranch
}

func (tree *declarationTree) addBranch(branch declarationBranch) {
	tree.branches = append(tree.branches, branch)
}

func (tree *declarationTree) grow(branch declarationBranch, keyName string, value interface{}, fieldLevel fieldLevelKind) {

	if keyName == "_package" {
		return
	}

	if isString(value) {
		branch.addSegment(keyName, valueKindString, fieldLevel)
		if branch.containsRecursiveness {
			tree.addBranch(branch)
			return
		}
		valueLiteral := fmt.Sprintf("%v", value)
		nextValue, isNotBasicType := tree.yamlData[valueLiteral]
		if !isNotBasicType {
			branch.addSegment(valueLiteral, valueKindString, fieldLevelZero)
			tree.addBranch(branch)
			return
		}

		tree.grow(branch, valueLiteral, nextValue, firstFieldLevel)
	}

	if isMap(value) {
		branch.addSegment(keyName, valueKindObject, fieldLevel)
		if branch.containsRecursiveness {
			tree.addBranch(branch)
			return
		}
		mapValue := value.(map[interface{}]interface{})
		for _key, _value := range mapValue {
			branchCopy := branch.copySelf()
			_keyName := fmt.Sprintf("%v", _key)
			tree.grow(branchCopy, _keyName, _value, fieldLevel+1)
		}
	}
}

func validateNonRecursiveness(yamlData map[interface{}]interface{}) (errs []error) {
	// var tree declarationTree
	// tree.grow()
	return
}
