package render

import "html/template"

func NewTemplate() (*template.Template, error) {
	return template.ParseGlob("./ui/templates/*html")
}
