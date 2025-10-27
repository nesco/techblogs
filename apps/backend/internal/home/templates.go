package home

import (
	_ "embed"
	"html/template"
)

//go:embed templates/index.html.tmpl
var homePageTemplateContent string

var homePageTemplate = template.Must(template.New("HomePage").Parse(homePageTemplateContent))
