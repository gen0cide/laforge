package competition

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net"
	"path/filepath"
	"time"

	"github.com/bradfitz/iter"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	entropySize = 32
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

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

func NewTemplate(tmpl string, includeScripts bool) *template.Template {

	tmplFuncs := template.FuncMap{
		"N":                   iter.N,
		"CustomIP":            CustomIP,
		"CustomInternalCNAME": CustomInternalCNAME,
		"CustomExternalCNAME": CustomExternalCNAME,
		"MyIP":                GetPublicIP,
	}

	tmp := template.New(RandomString(entropySize))

	if includeScripts {
		tmplFuncs["ScriptRender"] = ScriptRender
		tmplFuncs["DScript"] = DScript
	}

	tmp.Funcs(tmplFuncs)

	newTmpl, err := tmp.Parse(string(MustAsset(tmpl)))
	if err != nil {
		panic(err)
	}

	return newTmpl
}

func DScript(name string, c *Competition, e *Environment, i int, n *Network, h *Host, hn string) string {

	script := c.ParseScripts()[name]

	if script == nil {
		LogFatal("Script not found: " + name)
	}

	tmplFuncs := template.FuncMap{
		"N":                   iter.N,
		"CustomIP":            CustomIP,
		"CustomInternalCNAME": CustomInternalCNAME,
		"CustomExternalCNAME": CustomExternalCNAME,
		"MyIP":                GetPublicIP,
	}

	tmp := template.New(RandomString(entropySize))

	tmp.Funcs(tmplFuncs)
	newTmpl, err := tmp.Parse(string(script.Contents))
	if err != nil {
		panic(err)
	}

	filename := filepath.Join(e.TfScriptsDir(), fmt.Sprintf("%s%d_%s_%s", e.Prefix, i, hn, name))

	tb := TemplateBuilder{
		Competition: c,
		Environment: e,
		PodID:       i,
		Network:     n,
		Host:        h,
	}

	var tpl bytes.Buffer

	if err := newTmpl.Execute(&tpl, tb); err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, tpl.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
	return filename

}

func EmbedRender(t string, c *Competition, e *Environment, i int, n *Network, h *Host, hn string) string {
	tmpl := Render(t, c, e, i, n, h)
	filename := filepath.Join(e.TfScriptsDir(), fmt.Sprintf("%s%d_%s_%s", e.Prefix, i, hn, t))
	err := ioutil.WriteFile(filename, tmpl, 0644)
	if err != nil {
		panic(err)
	}
	return filename
}

func ScriptRender(t string, c *Competition, e *Environment, i int, n *Network, h *Host, hn string) string {
	if h == nil {
		h = &Host{
			Hostname: hn,
		}
	}
	tmpl := Render(t, c, e, i, n, h)
	filename := filepath.Join(e.TfScriptsDir(), fmt.Sprintf("%s%d_%s_%s", e.Prefix, i, hn, t))
	err := ioutil.WriteFile(filename, tmpl, 0644)
	if err != nil {
		panic(err)
	}
	return filename
}

func Render(tmpName string, c *Competition, e *Environment, t int, n *Network, h *Host) []byte {
	tb := TemplateBuilder{
		Competition: c,
		Environment: e,
		PodID:       t,
		Network:     n,
		Host:        h,
	}

	tmpl := NewTemplate(tmpName, false)

	var tpl bytes.Buffer

	if err := tmpl.Execute(&tpl, tb); err != nil {
		panic(err)
	}

	return tpl.Bytes()
}

func RenderTB(tmpName string, tb *TemplateBuilder) []byte {
	tmpl := NewTemplate(tmpName, true)

	var tpl bytes.Buffer

	if err := tmpl.Execute(&tpl, tb); err != nil {
		panic(err)
	}

	return tpl.Bytes()
}

func CustomInternalCNAME(e *Environment, n *Network, c string) string {
	return fmt.Sprintf("%s.%s.%s", c, n.Subdomain, e.Domain)
}

func CustomExternalCNAME(e *Environment, c string) string {
	return fmt.Sprintf("%s.%s", c, e.Domain)
}

func CustomIP(cidr string, offset, id int) string {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	newIP := ip.To4()
	lastOctet := offset + id
	newIP[3] = byte(lastOctet)
	return newIP.String()
}
