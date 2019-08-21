package model2markdown

type File struct {
	Path    string
	Package string
	Structs Structs
}

type Files []File
