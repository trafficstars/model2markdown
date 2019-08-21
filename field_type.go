package model2markdown

type FieldType string

func (t FieldType) StringForDocumentation() string {
	v := string(t)
	switch v {
	case `null.Int`:
		return `*int64`
	case `null.Uint`:
		return `*uint64`
	case `null.String`:
		return `*string`
	case `null.Boolean`:
		return `*bool`
	case `null.Float`:
		return `*float64`
	case `time.Time`:
		return `datetime`
	case `null.Time`:
		return `*datetime`
	default:
		return v
	}
}