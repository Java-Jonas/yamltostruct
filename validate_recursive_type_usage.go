package yamltostruct

func validateRecursiveTypeUsage(yamlData map[interface{}]interface{}) (errs []error) {
	tree := newDeclarationTree(yamlData)

	tree.grow(declarationBranch{}, "", yamlData, fieldLevelZero)

	for _, branch := range tree.branches {
		if branch.containsRecursiveness {
			errs = append(errs, newValidationErrorRecursiveTypeUsage(branch.declarationPath()))
		}
	}

	return
}
