package competition

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/bradfitz/iter"
)

const (
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	entropySize = 32
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	tmplFuncs  = template.FuncMap{
		"N":                   iter.N,
		"CustomIP":            CustomIP,
		"CustomInternalCNAME": CustomInternalCNAME,
		"CustomExternalCNAME": CustomExternalCNAME,
		"MyIP":                GetPublicIP,
		"GetUsersForHost":     GetUsersForHost,
		"GetUsersForOU":       GetUsersForOU,
		"GetAllUsers":         GetAllUsers,
		"Incr":                Incr,
		"SetZero":             SetZero,
		"CalculateReversePTR": CalculateReversePTR,
	}
)

type TemplateBuilder struct {
	Competition      *Competition
	Environment      *Environment
	PodID            int
	Network          *Network
	Host             *Host
	HostIndex        int
	ScriptErrorCount int
}

type TFTemplate struct {
	Name     string
	Builder  TemplateBuilder
	Template string
	Rendered string
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

func GetUsersForHost(c *Competition, h *Host) []User {
	users := []User{}
	for _, userGroup := range h.UserGroups {
		for _, u := range c.UserList[userGroup] {
			users = append(users, u)
		}
	}
	return users
}

func GetAllUsers(c *Competition) []User {
	users := []User{}
	for _, ug := range c.UserList {
		for _, u := range ug {
			users = append(users, u)
		}
	}
	return users
}

func GetUsersForOU(c *Competition, ou string) []User {
	return c.UserList[ou]
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}

func SetZero(val int) error {
	val = 0
	return nil
}

func NewTemplate(tmpl string, includeScripts bool) *template.Template {

	tmp := template.New(RandomString(entropySize))

	newFuncMap := template.FuncMap{}

	for funcName, funcHandler := range tmplFuncs {
		newFuncMap[funcName] = funcHandler
	}

	if includeScripts {
		newFuncMap["ScriptRender"] = ScriptRender
		newFuncMap["DScript"] = DScript
	}

	tmp.Funcs(newFuncMap)

	newTmpl, err := tmp.Parse(string(MustAsset(tmpl)))
	if err != nil {
		LogFatal("Error parsing template: template=" + tmpl + " error=" + err.Error())
	}

	return newTmpl
}

func DScript(name string, c *Competition, e *Environment, i int, n *Network, h *Host, hn string) string {

	script := c.ParseScripts()[name]

	if script == nil {
		LogFatal(fmt.Sprintf("Script not found! script=%s env=%s network=%s host=%s", name, e.Name, n.Name, h.Hostname))
	}

	tmp := template.New(RandomString(entropySize))

	tmp.Funcs(tmplFuncs)
	newTmpl, err := tmp.Parse(string(script.Contents))
	if err != nil {
		LogError("Error parsing script: script=" + name + " host=" + hn + " error=" + err.Error())
		return "SCRIPT_PARSING_ERROR"
	}

	filename := filepath.Join(e.TfScriptsDir(i), fmt.Sprintf("%s%d_%s_%s", e.Prefix, i, hn, name))

	tb := TemplateBuilder{
		Competition:      c,
		Environment:      e,
		PodID:            i,
		Network:          n,
		Host:             h,
		ScriptErrorCount: 0,
	}

	var tpl bytes.Buffer

	if err := newTmpl.Execute(&tpl, tb); err != nil {
		LogError("Error proccessing script: script=" + name + " host=" + hn + " error=" + err.Error())
		return "SCRIPT_PARSING_ERROR"
	}

	err = ioutil.WriteFile(filename, tpl.Bytes(), 0644)
	if err != nil {
		LogError("Error writing script: script=" + name + " path=" + filename)
		return "SCRIPT_PARSING_ERROR"
	}
	return filename

}

func EmbedRender(t string, c *Competition, e *Environment, i int, n *Network, h *Host, hn string) string {
	tmpl := Render(t, c, e, i, n, h)
	filename := filepath.Join(e.TfScriptsDir(i), fmt.Sprintf("%s%d_%s_%s", e.Prefix, i, hn, t))
	err := ioutil.WriteFile(filename, tmpl, 0644)
	if err != nil {
		LogError("Error writing embed script in EmbedRender(): path=" + filename)
		return filename
	}
	return filename
}

func StringRender(t string, c *Competition, e *Environment, i int, n *Network, h *Host) string {
	tb := TemplateBuilder{
		Competition:      c,
		Environment:      e,
		PodID:            i,
		Network:          n,
		Host:             h,
		ScriptErrorCount: 0,
	}

	tmp := template.New(RandomString(entropySize))

	tmp.Funcs(tmplFuncs)
	newTmpl, err := tmp.Parse(t)
	if err != nil {
		LogFatal("Error Parsing Template String in StringRender(): " + t)
	}

	var tpl bytes.Buffer

	if err := newTmpl.Execute(&tpl, tb); err != nil {
		LogFatal("Error Rendering Template String in StringRender(): " + t)
	}

	return tpl.String()
}

func ScriptRender(t string, c *Competition, e *Environment, i int, n *Network, h *Host, hn string) string {
	if h == nil {
		h = &Host{
			Hostname: hn,
		}
	}
	tmpl := Render(t, c, e, i, n, h)
	filename := filepath.Join(e.TfScriptsDir(i), fmt.Sprintf("%s%d_%s_%s", e.Prefix, i, hn, t))
	err := ioutil.WriteFile(filename, tmpl, 0644)
	if err != nil {
		LogError("Error writing script in ScriptRender(): path=" + filename)
		return filename
	}
	return filename
}

func Render(tmpName string, c *Competition, e *Environment, t int, n *Network, h *Host) []byte {
	tb := TemplateBuilder{
		Competition:      c,
		Environment:      e,
		PodID:            t,
		Network:          n,
		Host:             h,
		ScriptErrorCount: 0,
	}

	tmpl := NewTemplate(tmpName, false)

	var tpl bytes.Buffer

	if err := tmpl.Execute(&tpl, tb); err != nil {
		LogFatal("Error rendering template on Render(): " + tmpName)
	}

	return tpl.Bytes()
}

func RenderTB(tmpName string, tb *TemplateBuilder) []byte {
	tmpl := NewTemplate(tmpName, true)

	var tpl bytes.Buffer

	if err := tmpl.Execute(&tpl, tb); err != nil {
		LogFatal("Error rendering template on RenderTB(): " + tmpName + " error: " + err.Error())
	}

	return tpl.Bytes()
}

func Incr(val int) error {
	val = val + 1
	return nil
}

func RenderTBV2(tmpName string, tb *TemplateBuilder) []byte {
	tmpl := NewTemplate(tmpName, true)

	var tpl bytes.Buffer

	if err := tmpl.Execute(&tpl, tb); err != nil {
		LogFatal("Error rendering template on RenderTB(): " + tmpName + " error: " + err.Error())
	}

	return tpl.Bytes()
}

func NewTemplateContext(c *Competition, e *Environment, pid int, n *Network, h *Host) *TemplateBuilder {
	return &TemplateBuilder{
		Competition:      c,
		Environment:      e,
		PodID:            pid,
		Network:          n,
		Host:             h,
		ScriptErrorCount: 0,
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
