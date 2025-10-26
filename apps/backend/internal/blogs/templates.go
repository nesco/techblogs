package blogs

import (
	_ "embed"
	"html/template"
)

//go:embed templates/blog_list.html.tmpl
var blogListTemplateContent string

var BlogListTemplate = template.Must(template.New("BlogList").Parse(blogListTemplateContent))

//go:embed templates/blog_feed.xml.tmpl
var blogFeedTemplateContent string

var BlogFeedTemplate = template.Must(template.New("BlogFeed").Parse(blogFeedTemplateContent))
