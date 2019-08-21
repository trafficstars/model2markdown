package templates

import (
	"strings"
	"text/template"
)

const (
	modelMarkdownTemplateText = ``+
`# [{{ .File.Path }}] {{ .Struct.Name }}
| name | sql{{- if .shouldPrintJSON}} | json{{- end }} | description |
| ---- | ---{{- if .shouldPrintJSON}} | ----{{- end }} | ----------- |
{{- $shouldPrintJSON := .shouldPrintJSON }}
{{- range $index,$field := .Struct.Fields }}
| {{ escapeMarkdown $field.Name }} | {{ escapeMarkdown $field.SQLFieldName }}{{- if $shouldPrintJSON}} | {{ escapeMarkdown $field.JSONFieldName }}{{- end }} | {{ escapeMarkdown ( stringsJoin $field.Comments "," ) }} |
{{- end }}
`
)

var (
	ModelMarkdownTemplate *template.Template
)

func init() {
	var err error
	tpl := template.New(`ModelMarkdownTemplate`).Funcs(map[string]interface{}{
		"escapeMarkdown" : func(in string) string {
			// Implementation of this function is copied from
			// https://raw.githubusercontent.com/ekalinin/github-markdown-toc.go/09cbee650f0f3f0974d2959467713637dbd99f41/internals.go
			specChar := []string{"\\", "`", "*", "_", "{", "}", "#", "+", "-", ".", "!"}
			res := in

			for _, c := range specChar {
				res = strings.Replace(res, c, "\\"+c, -1)
			}
			return res
		},
		"stringsJoin": func(slice []string, sep string) string {
			return strings.Join(slice, sep)
		},
	})

	ModelMarkdownTemplate, err = tpl.Parse(modelMarkdownTemplateText)
	if err != nil {
		panic(err)
	}
}