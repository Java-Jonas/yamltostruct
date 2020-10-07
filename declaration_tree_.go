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
	// current level is the object itself
	fieldLevelZero fieldLevelKind = iota
	// current level is the root level of the object
	firstFieldLevel
	secondFieldLevel
)

type declarationSegment struct {
	keyName   string
	valueKind valueKind
	// TODO typeKind value/reference
	fieldLevel fieldLevelKind
}

type declarationBranch struct {
	segments []declarationSegment
	// TODO: branchClosureKind (recursiveness, reference type ..)
	containsRecursiveness bool
}

func (branch *declarationBranch) addSegment(keyName string, valueKind valueKind, fieldLevel fieldLevelKind) {
	// TODO: maybe this shouldn't be done here
	for _, segment := range branch.segments {
		if segment.keyName == keyName {
			branch.containsRecursiveness = true
		}
	}
	branch.segments = append(branch.segments, declarationSegment{keyName, valueKind, fieldLevel})
}

// we list the segments' typeNames with some additional logic
// to be more explicit (paths to nested field names will
// be concatenated eg. "foo.bar")
func (branch declarationBranch) declarationPath() []string {
	var path []string

	var wasStructField bool
	var parentStructName string

	for _, segment := range branch.segments {
		if segment.valueKind == valueKindObject {
			wasStructField = true
			parentStructName = segment.keyName
			continue
		}

		if wasStructField && segment.fieldLevel == secondFieldLevel {
			path = append(path, parentStructName+"."+segment.keyName)
		} else {
			path = append(path, segment.keyName)
		}
		wasStructField = false
		parentStructName = ""

	}

	if wasStructField {
		path = append(path, parentStructName)
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

func newDeclarationTree(yamlData map[interface{}]interface{}) *declarationTree {
	return &declarationTree{yamlData: yamlData}
}

func (tree *declarationTree) addBranch(branch declarationBranch) {
	tree.branches = append(tree.branches, branch)
}

// a recursive function to travel through the yamlData, creating
// a different branch for each path
func (tree *declarationTree) grow(branch declarationBranch, keyName string, value interface{}, fieldLevel fieldLevelKind) {

	if fieldLevel == firstFieldLevel && keyName == "_package" {
		return
	}

	if isString(value) {
		branch.addSegment(keyName, valueKindString, fieldLevel)
		if branch.containsRecursiveness {
			// detected recursiveness implies this is the end of the branch
			tree.addBranch(branch)
			return
		}
		valueLiteral := fmt.Sprintf("%v", value)
		// if a key cannot be found at root level we assume it's a basic type eg. string, int ..
		nextValue, isNotBasicType := tree.yamlData[valueLiteral]

		/*
			if isReferencingValue(){
				branch.addSegment(valueLiteral, valueKindString, fieldLevelZero)
				// a reference type implies this is the end of the branch
				tree.addBranch(branch)
				return
			}

		*/

		if !isNotBasicType {
			branch.addSegment(valueLiteral, valueKindString, fieldLevelZero)
			// a basic type implies this is the end of the branch
			tree.addBranch(branch)
			return
		}
		// we get here only when the value is a reference to a user defined type
		tree.grow(branch, valueLiteral, nextValue, firstFieldLevel)
	}

	if isMap(value) {
		branch.addSegment(keyName, valueKindObject, fieldLevel)
		if branch.containsRecursiveness {
			// detected recursiveness implies this is the end of the branch
			tree.addBranch(branch)
			return
		}
		mapValue := value.(map[interface{}]interface{})
		for _key, _value := range mapValue {
			// the branch is copied; this is basically a fork
			branchCopy := branch.copySelf()
			_keyName := fmt.Sprintf("%v", _key)
			// we go a level deeper (fieldLevel+1) and handle each key/value
			// pair in the next tree.grow() execution
			tree.grow(branchCopy, _keyName, _value, fieldLevel+1)
		}
	}
}
