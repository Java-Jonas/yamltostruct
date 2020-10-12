package yamltostruct

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// func extractMapKeyRecursive(mapKeys []string, valueString string) []string {
// 	re := regexp.MustCompile(`map\[(.*)\]`)
// 	match := re.FindString(valueString)
// 	if match == "" {
// 		return mapKeys
// 	}
// 	// matchLocation := re.FindIndex([]byte(valueString))
// 	// subsequentValueString := valueString[matchLocation[1]:]
// 	mapKey := match[4 : len(match)-1]
// 	mapKeys = append(mapKeys, mapKey)
// 	mapKeys = extractMapKeyRecursive(mapKeys, mapKey)
// 	// mapKeys = extractMapKeyRecursive(mapKeys, subsequentValueString)
// 	return mapKeys
// }

func extracpMapDeclExpression(file *ast.File) ast.Expr {
	return file.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type
}

func extractMapKeyRecursive(mapKeys []string, typeSpec ast.Expr) []string {
	return []string{}
}

func extractMapKeys(valueString string) []string {
	mockSrc := `
	package main
	type mockType` + valueString

	file, _ := parser.ParseFile(token.NewFileSet(), "", mockSrc, 0)
	return extractMapKeyRecursive([]string{}, extracpMapDeclExpression(file))
}

func extractRootLevelMapDeclarations(valueString string) []string {
	return []string{}
}
