package blogs

import (
	_ "embed"
	"html/template"
	texttemplate "text/template"
)

//go:embed templates/blog_list.html.tmpl
var blogListTemplateContent string

var BlogListTemplate = template.Must(template.New("BlogList").Parse(blogListTemplateContent))

//go:embed templates/blog_feed.xml.tmpl
var blogFeedTemplateContent string

var BlogFeedTemplate = texttemplate.Must(texttemplate.New("BlogFeed").Parse(blogFeedTemplateContent))
