package templates

import (
	"embed"
	"text/template"
	"strings"
)

//go:embed *.tmpl
var FS embed.FS

func ParseTemplate(templateName string,data any) (string, error){
	tmpl := template.New("")

	parsed, err := tmpl.ParseFS(FS,"*.tmpl")
	if err != nil{
		return "",err
	}

	var builder strings.Builder
	if err := parsed.ExecuteTemplate(&builder,templateName,data); err != nil{
		//fmt.Errorf("Loading template: %w",err)
		return "", err
	}
	return builder.String(),nil
}
