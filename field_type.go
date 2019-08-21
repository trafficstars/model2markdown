package model2markdown

type FieldType string

func (t FieldType) StringForDocumentation() string {
	v := string(t)
	switch v {
	case `Int32Array`:
		return `[]int32`
	case `Int64Array`:
		return `[]int64`
	case `StringArray`:
		return `[]string`
	case `Int`:
		return `*int64`
	case `Uint`:
		return `*uint64`
	case `String`:
		return `*string`
	case `Boolean`, `Bool`:
		return `*bool`
	case `Float`:
		return `*float64`
	case `Time`:
		return `datetime`
	case `*Time`:
		return `*datetime`
	default:
		return v
	}
}
