package competition

type ScriptMap map[string]*Script

type Script struct {
	Name string
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
