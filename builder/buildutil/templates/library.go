package templates

import (
	"bytes"
	"sync"
	"text/template"

	"github.com/gen0cide/laforge/core"

	"github.com/pkg/errors"
)

// Book is a reference to a specific template
type Book struct {
	Name     string
	Original []byte
	Template *template.Template
}

// Library is a holder for all the templates used by a builder
type Library struct {
	sync.RWMutex
	Books map[string]*Book
}

// NewLibrary returns a new zero value Library
func NewLibrary() *Library {
	return &Library{
		Books: map[string]*Book{},
	}
}

// AddBook adds a new template to the library index
func (l *Library) AddBook(name string, data []byte) (*Book, error) {
	l.Lock()
	defer l.Unlock()
	if b, ok := l.Books[name]; ok {
		return b, nil
	}

	t := template.New(name)
	t.Funcs(core.TemplateFuncLib)
	newT, err := t.Parse(string(data))
	if err != nil {
		return nil, err
	}

	b := &Book{
		Name:     name,
		Original: data,
		Template: newT,
	}

	l.Books[name] = b
	return b, nil
}

// Execute uses the denoted book and renders a template based off of the passed context
func (l *Library) Execute(name string, context *Context) ([]byte, error) {
	l.Lock()
	defer l.Unlock()

	book, found := l.Books[name]
	if !found {
		return []byte{}, errors.Errorf("could not locate template book named %s", name)
	}

	buf := new(bytes.Buffer)
	err := book.Template.Execute(buf, context)
	return buf.Bytes(), err
}

// ExecuteGroup uses the denoted book with all secondary books chained and renders a template based off of the passed context
func (l *Library) ExecuteGroup(baseID string, appendIDs []string, context *Context) ([]byte, error) {
	l.Lock()
	defer l.Unlock()

	book, found := l.Books[baseID]
	if !found {
		return []byte{}, errors.Errorf("could not locate template book named %s", baseID)
	}
	appendMap := map[string]bool{}
	for _, x := range appendIDs {
		if x == baseID {
			continue
		}
		appendMap[x] = true
	}
	tmpl := book.Template
	var err error
	for bn, bk := range l.Books {
		_, ok := appendMap[bn]
		if !ok {
			continue
		}
		tmpl, err = tmpl.Parse(string(bk.Original))
		if err != nil {
			panic(err)
		}
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, context)
	return buf.Bytes(), err
}
