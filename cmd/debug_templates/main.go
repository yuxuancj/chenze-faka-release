package main

import (
	"chenze-faka/web"
	"fmt"
	"text/template"
)

func main() {
	tmpl, err := template.ParseFS(web.StaticFiles, "templates/**/*.html")
	if err != nil {
		panic(err)
	}
	fmt.Println("Root name:", tmpl.Name())
	for _, t := range tmpl.Templates() {
		fmt.Printf("  Template: name=%q\n", t.Name())
	}
}
