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
		for _, sqlTagName := range []string{`sql`,`gorm`,`reform`} {
			sqlFieldName = strings.Split(tag.Get(sqlTagName), `,`)[0]
			if sqlFieldName != `` {
				break
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

		r = append(r, Field{
			Name: name,
			Type: FieldType(f.Type.Pos()),
			SQLFieldName: sqlFieldName,
			JSONFieldName:jsonFieldName,
			Comments:comments,
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
				Name: name,
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
