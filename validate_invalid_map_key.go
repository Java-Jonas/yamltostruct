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

func findMapType(expr ast.Expr) *ast.MapType {
	mapType, ok := expr.(*ast.MapType)
	if ok {
		return mapType
	}
	arrayType, ok := expr.(*ast.ArrayType)
	if ok {
		return findMapType(arrayType.Elt)
	}
	starType, ok := expr.(*ast.StarExpr)
	if ok {
		return findMapType(starType.X)
	}
	return nil
}

func extractMapKeyRecursive(mapKeys []string, mapType *ast.MapType, mockSrc string) []string {

	if mapType == nil || mapType.Key == nil {
		return mapKeys
	}

	keyIdent, ok := mapType.Key.(*ast.Ident)
	if !ok {
		mapKey := mockSrc[mapType.Key.Pos()-1 : mapType.Key.End()-1]
		mapKeys = append(mapKeys, mapKey)
	} else {
		mapKeys = append(mapKeys, keyIdent.Name)
	}

	mapKeys = extractMapKeyRecursive(mapKeys, findMapType(mapType.Value), mockSrc)

	return mapKeys
}

func extractMapKeys(valueString string) []string {
	mockSrc := `
	package main
	type mockType ` + valueString

	file, _ := parser.ParseFile(token.NewFileSet(), "", mockSrc, 0)

	typeExpression := extracpMapDeclExpression(file)
	mapType, ok := typeExpression.(*ast.MapType)
	if !ok {
		return []string{}
	}

	return extractMapKeyRecursive([]string{}, mapType, mockSrc)
}

func extractRootLevelMapDeclarations(valueString string) []string {
	return []string{}
}
