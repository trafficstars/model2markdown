package model2markdown

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/xaionaro-go/errors"
)

func exprToTypeName(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.StarExpr:
		return `*` + exprToTypeName(v.X)
	case *ast.SelectorExpr:
		return v.Sel.Name
	case *ast.ArrayType:
		if v.Len == nil {
			return fmt.Sprintf(`[]%v`, exprToTypeName(v.Elt))
		}
		return fmt.Sprintf(`[%v]%v`, v.Len, exprToTypeName(v.Elt))
	case *ast.MapType:
		return fmt.Sprintf(`map[%v]%v`, exprToTypeName(v.Key), exprToTypeName(v.Value))
	case *ast.InterfaceType:
		return fmt.Sprintf(`%v`, v.Interface)
	}

	panic(fmt.Sprintf(`unknown type: %T`, expr))
}

func fieldGetStructTag(f *ast.Field) reflect.StructTag {
	if f.Tag != nil {
		tag := f.Tag.Value
		if len(tag) >= 3 {
			return reflect.StructTag(strings.Trim(tag, `"`))
		}
	}

	return reflect.StructTag("")
}

func structTypeGetFields(structType *ast.StructType) (r Fields) {
	for _, f := range structType.Fields.List {
		if len(f.Names) == 0 {
			continue
		}

		if len(f.Names) > 1 {
			panic(fmt.Errorf("not implemented case: %v", f.Names))
		}
		name := f.Names[0].Name

		tag := fieldGetStructTag(f)

		var sqlFieldName string
		for _, sqlTagName := range []string{`sql`, `reform`} {
			sqlFieldName = strings.Split(tag.Get(sqlTagName), `,`)[0]
			if sqlFieldName != `` {
				break
			}
		}
		if sqlFieldName == `` {
			for _, gormTag := range strings.Split(tag.Get(`gorm`), `,`) {
				parts := strings.Split(gormTag, `:`)
				if len(parts) != 2 {
					continue
				}
				key := parts[0]
				value := parts[1]
				switch key {
				case `column`:
					sqlFieldName = value
				}
			}
		}
		if sqlFieldName == `` {
			sqlFieldName = gorm.ToColumnName(name)
		}

		var jsonFieldName string
		if jsonFieldName == `` {
			jsonFieldName = strings.Split(tag.Get(`json`), `,`)[0]
		}
		if jsonFieldName == `` {
			jsonFieldName = name
		}

		var comments []string
		if f.Comment != nil {
			for _, comment := range f.Comment.List {
				comments = append(comments, strings.Trim(comment.Text, `/* \t\n\r`))
			}
		}

		typeName := exprToTypeName(f.Type)

		r = append(r, Field{
			Name:          name,
			Type:          FieldType(typeName),
			SQLFieldName:  sqlFieldName,
			JSONFieldName: jsonFieldName,
			Comments:      comments,
		})
	}

	return
}

func declsGetStructs(decls []ast.Decl) (r Structs) {
	for _, decl := range decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			if structType.Incomplete {
				continue
			}

			name := typeSpec.Name.Name
			if name == "" {
				continue
			}

			r = append(r, Struct{
				Name:   name,
				Fields: structTypeGetFields(structType),
			})
		}
	}
	return
}

func ParseFile(filePath string) (r File, err error) {
	parsed, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.ParseComments)
	if err != nil {
		err = errors.Wrap(err, `unable to parse the file`, filePath)
		return
	}

	if parsed.Name == nil {
		err = fmt.Errorf(`invalid golang file %v: package name not found`, filePath)
		return
	}

	r.Package = parsed.Name.Name

	r.Structs = declsGetStructs(parsed.Decls)

	r.Path = filePath

	return
}
