package competition

import (
	"errors"
	"io/ioutil"
)

type ScriptMap map[string]*Script

type Script struct {
	Contents []byte
}

func LoadScript(path string) (*Script, error) {
	script := Script{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		LogError("Could not read file: " + path)
		return nil, errors.New("invalid script file")
	}
	script.Contents = data
	return &script, nil
}

func (s *Script) RenderString(t TemplateBuilder) string {
	return ""
}

func (s *Script) RenderFile(path string, t TemplateBuilder) error {
	return nil
}

func GetScripts() ScriptMap {
	return make(map[string]*Script)
}
