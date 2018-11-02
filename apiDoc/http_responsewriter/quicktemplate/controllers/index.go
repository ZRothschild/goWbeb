package controllers

import (
	"../templates"
	"github.com/kataras/iris"
)
// Index renders our ../templates/index.qtpl file using the compiled ../templates/index.qtpl.go file.
func Index(ctx iris.Context) {
	tmpl := &templates.Index{}
	// render the template
	ExecuteTemplate(ctx, tmpl)
}