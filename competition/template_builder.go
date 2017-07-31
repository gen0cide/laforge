package competition

import (
	"bytes"
	"errors"
	"html/template"
	"path/filepath"
	"reflect"
	"strings"
)

type TemplateBuilder struct {
	Competition *Competition
	Environment *Environment
	Pod         *Pod
	Network     *Network
	Host        *Host
}

type TFTemplate struct {
	Name     string
	Builder  TemplateBuilder
	Template string
	Rendered string
}

func TFRender(tfobj interface{}) (string, error) {
	tmplName := strings.ToLower(reflect.TypeOf(tfobj).Name())
	filename := tmplName + ".tmpl"
	tmplFile := filepath.Join("terraform", filename)
	asset, err := Asset(tmplFile)
	if err != nil {
		return "", err
	}
	var tmplBuffer bytes.Buffer
	tmpl, err := template.New(tmplName).Parse(string(asset))
	if err != nil {
		return "", errors.New("Fatal Error parsing terraform template (" + tmplFile + "): " + err.Error())
	}
	if err := tmpl.Execute(&tmplBuffer, tfobj); err != nil {
		return "", errors.New("Fatal Error rendering terraform template (" + tmplFile + "): " + err.Error())
	}

	return tmplBuffer.String(), nil
}
