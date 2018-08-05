package laforge

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/hashicorp/hcl/hcl/printer"

	"github.com/gen0cide/laforge/static"
)

// RenderHCLv2Object is a generic templating function for HCLv2 compatible types
func RenderHCLv2Object(i interface{}) ([]byte, error) {
	t := reflect.TypeOf(i)
	tname := strings.ToLower(t.Name())
	tmplname := fmt.Sprintf("%s.hcl.tmpl", tname)
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
