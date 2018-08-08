package core

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/hashicorp/hcl/hcl/printer"

	"github.com/gen0cide/laforge/static"
)

// RenderHCLv2Object is a generic templating function for HCLv2 compatible types
func RenderHCLv2Object(i interface{}) ([]byte, error) {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		v := reflect.ValueOf(i)
		t = reflect.Indirect(v).Type()
	}
	tname := strings.ToLower(strcase.ToSnake(t.Name()))
	tmplname := fmt.Sprintf("%s.laforge.tmpl", tname)
	Logger.Debugf("Searching for template %s", tmplname)
	tmpldata, err := static.ReadFile(tmplname)
	if err != nil {
		return []byte{}, err
	}
	tm, err := template.New(tmplname).Parse(string(tmpldata))
	if err != nil {
		return []byte{}, err
	}
	buf := new(bytes.Buffer)
	err = tm.Execute(buf, i)
	if err != nil {
		return []byte{}, err
	}
	return printer.Format(buf.Bytes())
}
