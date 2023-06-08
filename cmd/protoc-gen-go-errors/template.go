package main

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"
)

const errorsTpl = `
{{ range .Errors }}
var {{.LowerCamelValue}} *eerrors.GoError
{{- end }}

var i18n = map[string]map[string]string{
{{- range .Errors }}
	"{{.Key}}": map[string]string{
		{{- range $k,$v :=  .I18n }}
			"{{$k}}": "{{$v}}",
		{{- end }}
	},
{{- end }}
}

func ReasonI18n(e eerrors.Error, lan string) string {
	return i18n[eerrors.FromError(e).Reason][lan]
}

func init() {
{{- range .Errors }}
{{.LowerCamelValue}} = eerrors.New(int(codes.{{.Code}}), "{{.Key}}", {{.Name}}_{{.Value}}.String())
eerrors.Register({{.LowerCamelValue}})
{{- end }}
}

{{ range .Errors }}
{{if .HasComment}}{{.Comment}}{{end}}func {{.UpperCamelValue}}() eerrors.Error {
	 return {{.LowerCamelValue}}
}
{{ end }}
`

type errorInfo struct {
	Name            string
	Value           string
	Code            string
	UpperCamelValue string
	LowerCamelValue string
	Key             string
	Comment         string
	HasComment      bool
	I18n            map[string]string
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTpl)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return string(GoFmt(buf.Bytes()))
}

// GoFmt 格式化代码
func GoFmt(buf []byte) []byte {
	formatted, err := format.Source(buf)
	if err != nil {
		panic(fmt.Errorf("%s\nOriginal code:\n%s", err.Error(), buf))
	}
	return formatted
}
