package model2markdown

type Field struct {
	Name          string
	Type          FieldType
	JSONFieldName string
	SQLFieldName  string
	Comments      []string
}

type Fields []Field
