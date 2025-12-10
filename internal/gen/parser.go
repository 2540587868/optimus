package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
)

type ParseConfig struct {
	FileName   string
	StructName string
}

type ParseResult struct {
	Imports []string
	Fields  []FieldInfo
}

func ParseFile(cfg ParseConfig) (*ParseResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, cfg.FileName, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	srcContent, err := os.ReadFile(cfg.FileName)
	if err != nil {
		return nil, err
	}

	result := &ParseResult{}

	for _, imp := range node.Imports {
		var sb strings.Builder
		if imp.Name != nil {
			sb.WriteString(imp.Name.Name)
			sb.WriteString(" ")
		}
		sb.WriteString(imp.Path.Value)
		result.Imports = append(result.Imports, sb.String())
	}

	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != cfg.StructName {
			return true
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		found = true
		for _, field := range st.Fields.List {
			if len(field.Names) == 0 {
				continue
			}

			info := parseField(field, fset, srcContent)
			if info != nil {
				result.Fields = append(result.Fields, *info)
			}
		}
		return false
	})

	if !found {
		return nil, fmt.Errorf("struct %s not found in %s", cfg.StructName, cfg.FileName)
	}

	return result, nil
}

func parseField(field *ast.Field, fset *token.FileSet, src []byte) *FieldInfo {
	name := field.Names[0].Name

	tagVal := ""
	defaultTag := ""
	if field.Tag != nil {
		cleanTag := strings.Trim(field.Tag.Value, "`")
		parsedTag := reflect.StructTag(cleanTag)
		tagVal = parsedTag.Get("opt")
		defaultTag = parsedTag.Get("default")
	}

	if tagVal == "-" {
		return nil
	}

	typeStart := fset.Position(field.Type.Pos()).Offset
	typeEnd := fset.Position(field.Type.End()).Offset
	fieldType := string(src[typeStart:typeEnd])

	defaultValCode := ""
	if defaultTag != "" {
		if fieldType == "string" {
			defaultValCode = fmt.Sprintf("%q", defaultTag)
		} else {
			defaultValCode = defaultTag
		}
	}

	var comments []string
	if field.Doc != nil {
		for _, c := range field.Doc.List {
			comments = append(comments, c.Text)
		}
	}

	isSlice := false
	isPointer := false

	isMap := false
	mapKType := ""
	mapVType := ""

	paramType := fieldType
	funcName := "With" + name

	if arrayType, ok := field.Type.(*ast.ArrayType); ok {
		isSlice = true
		funcName = "Add" + name

		eltStart := fset.Position(arrayType.Elt.Pos()).Offset
		eltEnd := fset.Position(arrayType.Elt.End()).Offset
		paramType = string(src[eltStart:eltEnd])
	} else if startExpr, ok := field.Type.(*ast.StarExpr); ok {
		innerStart := fset.Position(startExpr.X.Pos()).Offset
		innerEnd := fset.Position(startExpr.X.End()).Offset
		innerType := string(src[innerStart:innerEnd])
		if IsBasicType(innerType) {
			isPointer = true
			paramType = innerType
		} else {
			isPointer = false
			paramType = fieldType
		}
	} else if mapType, ok := field.Type.(*ast.MapType); ok {
		isMap = true
		funcName = "Add" + name

		kStart := fset.Position(mapType.Key.Pos()).Offset
		kEnd := fset.Position(mapType.Key.End()).Offset
		mapKType = string(src[kStart:kEnd])

		vStart := fset.Position(mapType.Value.Pos()).Offset
		vEnd := fset.Position(mapType.Value.End()).Offset
		mapVType = string(src[vStart:vEnd])
	} else {
		paramType = fieldType
	}

	if tagVal != "" {
		funcName = tagVal
	}

	return &FieldInfo{
		Name:       name,
		Type:       fieldType,
		Func:       funcName,
		ParamType:  paramType,
		IsSlice:    isSlice,
		IsPointer:  isPointer,
		Comments:   comments,
		DefaultVal: defaultValCode,
		IsMap:      isMap,
		KeyType:    mapKType,
		ValueType:  mapVType,
	}
}

func IsBasicType(typeName string) bool {
	switch typeName {
	case "bool", "string",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"uintptr", "byte", "rune",
		"float32", "float64",
		"complex64", "complex128",
		"time.Duration":
		return true
	default:
		return false
	}
}
