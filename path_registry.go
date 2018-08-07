package laforge

// PathRegistry is a type that tracks the relative file paths of state configurations that include external sources
type PathRegistry struct {
	DB map[CallFile]*PathResolver // key = composite of $type.$id, value = PathResolver
}

// PathResolver defines the mapping of paths declared in a CallFile and their mapping to files on local disk or lack of resolution
type PathResolver struct {
	Mapping    map[string]*LocalFileRef // map[provided_path_declaration] => LocalFileRef
	Unresolved map[string]bool          // map[provided_path_declaration] => true when we can't resolve it
}

// LocalFileRef is a basic type to hold information about a resolved file that was declared inside a state declaration
type LocalFileRef struct {
	Base          string
	AbsPath       string
	RelPath       string
	Cwd           string
	DeclaredPath  string
	RelToCallFile string
}
