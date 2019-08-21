package model2markdown

import (
	"github.com/trafficstars/model2markdown/templates"
	"os"
	"path/filepath"
)

type Library struct {
	Files Files
}

func NewLibrary() *Library {
	return &Library{}
}

func (lib *Library) AddFile(f File) {
	lib.Files = append(lib.Files, f)
}

type libraryStruct struct {
	Struct
	File *File
}

func (s *libraryStruct) GetDocumentFileName() string {
	return s.File.Package+`_`+s.Struct.Name+`.md`
}

func (lib *Library) GenerateMarkdownsToDirectory(outputDir string) error {
	structs := map[string]*libraryStruct{}

	for idx := range lib.Files {
		file := &lib.Files[idx]
		for _, oneStruct := range file.Structs {
			s := &libraryStruct{
				oneStruct,
				file,
			}
			fileName := s.GetDocumentFileName()
			if structs[fileName] != nil {
				return NewErrConflict(*s, *structs[fileName])
			}
			structs[fileName] = s
		}
	}

	for fileName, s := range structs {
		file, err := os.OpenFile(filepath.Join(outputDir, fileName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		shouldPrintJSON := false
		for _, field := range s.Fields {
			if field.Name != field.JSONFieldName {
				shouldPrintJSON = true
				break
			}
		}

		err = templates.ModelMarkdownTemplate.Execute(file, map[string]interface{}{
			"File": s.File,
			"Struct": s.Struct,
			"shouldPrintJSON": shouldPrintJSON,
		})
		if err != nil {
			return err
		}
	}

	return nil
}