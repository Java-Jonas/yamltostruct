package yamltostruct

import (
	"fmt"
)

type declarationBranch []string
type declarationTree struct {
	yamlData map[interface{}]interface{}
	branches []declarationBranch
}

func (db declarationBranch) copySelf() declarationBranch {
	branchCopy := make(declarationBranch, len(db))
	copy(branchCopy, db)
	return branchCopy
}

func (db declarationBranch) wouldCreateLoopWith(keyName string) (wouldCreateLoop bool) {

	for _, segment := range db {
		if segment == keyName {
			return true
		}
	}

	return false
}

func (tree *declarationTree) addBranch(branch declarationBranch) {
	tree.branches = append(tree.branches, branch)
}

func (tree *declarationTree) grow(branch declarationBranch, keyName string, value interface{}) {

	if isString(value) {
		if branch.wouldCreateLoopWith(keyName) {
			branch = append(branch, keyName)
			tree.addBranch(branch)
			return
		}
		branch = append(branch, keyName)
		valueLiteral := fmt.Sprintf("%v", value)
		nextValue, ok := tree.yamlData[valueLiteral]
		if !ok {
			branch = append(branch, valueLiteral)
			tree.addBranch(branch)
			return
		}

		tree.grow(branch, valueLiteral, nextValue)
	}

	if isMap(value) {
		branch = append(branch, keyName+".")
		mapValue := value.(map[interface{}]interface{})
		for _key, _value := range mapValue {
			branchCopy := branch.copySelf()
			_keyName := fmt.Sprintf("%v", _key)
			tree.grow(branchCopy, _keyName, _value)
		}
	}
}

func validateNonRecursiveness(yamlData map[interface{}]interface{}) (errs []error) {
	return
}
