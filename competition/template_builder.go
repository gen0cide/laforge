package competition

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"reflect"
	"strings"
)

type TemplateBuilder struct {
	Competition *Competition
	Environment *Environment
	PodID       int
	Network     *Network
	Host        *Host
	HostIndex   int
}

func (t *TemplateBuilder) EnvItemName(i interface{}, opts ...string) string {
	optString := ""
	if len(opts) > 0 {
		optString = fmt.Sprintf("_%s", strings.Join(opts, "_"))
	}
	name := strings.ToLower(reflect.TypeOf(i).Name())
	return fmt.Sprintf("%s_%s%s", t.Environment.Suffix(t.PodID), name, optString)
}

func (t *TemplateBuilder) NetItemName(i interface{}, opts ...string) string {
	optString := ""
	if len(opts) > 0 {
		optString = fmt.Sprintf("_%s", strings.Join(opts, "_"))
	}
	name := strings.ToLower(reflect.TypeOf(i).Name())
	return fmt.Sprintf("%s_%s_%s%s", t.Environment.Suffix(t.PodID), t.Network.Name, name, optString)
}

type TFTemplate struct {
	Name     string
	Builder  TemplateBuilder
	Template string
	Rendered string
}

func NewTemplateContext(c *Competition, e *Environment, pid int, n *Network, h *Host) *TemplateBuilder {
	return &TemplateBuilder{
		Competition: c,
		Environment: e,
		PodID:       pid,
		Network:     n,
		Host:        h,
	}
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
