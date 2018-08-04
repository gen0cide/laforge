package competition

type FileMap map[string]*File

type File struct {
	Name string
}

func GetFiles() FileMap {
	return make(map[string]*File)
}
