package templates

import "text/template"

type Book struct {
	Name     string
	Filename string
	Filepath string
	Original []byte
}

type Library struct {
	Books map[string]*template.Template
}
