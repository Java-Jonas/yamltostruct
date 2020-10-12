package yamltostruct

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func extracpMapDeclExpression(file *ast.File) ast.Expr {
	return file.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type
}

func findMapTypeRecursive(expr ast.Expr) *ast.MapType {
	if mapType, ok := expr.(*ast.MapType); ok {
		return mapType
	}
	if arrayType, ok := expr.(*ast.ArrayType); ok {
		return findMapTypeRecursive(arrayType.Elt)
	}
	if starType, ok := expr.(*ast.StarExpr); ok {
		return findMapTypeRecursive(starType.X)
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

	mapKeys = extractMapKeyRecursive(mapKeys, findMapTypeRecursive(mapType.Value), mockSrc)

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

func validateMapKeys(valueString string, yamlData map[interface{}]interface{}) []string {
	mapKeys := extractMapKeys(valueString)
	if len(mapKeys) == 0 {
		return nil
	}

	return []string{}
}
