package laforge

import "github.com/imdario/mergo"

// Laforge defines the type that holds the global namespace within the laforge configuration engine
type Laforge struct {
	Filename    string               `json:"filename"`
	Includes    []string             `json:"include,omitempty"`
	BaseDir     string               `hcl:"base_dir,attr" json:"base_dir,omitempty"`
	CurrDir     string               `json:"current_dir,omitempty"`
	User        User                 `hcl:"user,block" json:"user,omitempty"`
	Competition Competition          `hcl:"competition,block" json:"competition,omitempty"`
	Environment Environment          `hcl:"environment,block" json:"environment,omitempty"`
	Hosts       []*Host              `hcl:"host,block" json:"hosts,omitempty"`
	Networks    []*Network           `hcl:"network,block" json:"networks,omitempty"`
	Identities  []*Identity          `hcl:"identity,block" json:"identities,omitempty"`
	Scripts     []*Script            `hcl:"script,block" json:"scripts,omitempty"`
	Commands    []*Command           `hcl:"command,block" json:"commands,omitempty"`
	HostMap     map[string]*Host     `json:"-"`
	NetworkMap  map[string]*Network  `json:"-"`
	IdentityMap map[string]*Identity `json:"-"`
	ScriptMap   map[string]*Script   `json:"-"`
	CommandMap  map[string]*Command  `json:"-"`
	Caller      Caller               `json:"-"`
}

// Opt defines a basic HCLv2 option label:
//	config "keyName" {
//		value = "valueData"
//	}
type Opt struct {
	Key   string `hcl:",label" json:"key,omitempty"`
	Value string `hcl:"value,attr" json:"value,omitempty"`
}

// OnConflict defines a configuration override for how to handle conflicting objects
// in the differential enumeration of objects. For example, if you have a file that has loaded
// a host definition (host "foo" {...}) and you have included an additional configuration
// in your env.laforge, you can specify an on_conflict { do = "" } setting within
// your block to control how the state of host "foo" has state changes applied to it.
//	host "foo" {
//		on_conflict {
//			do = "overwrite"
//		}
//	}
// The default behavior (no on_conflict block) is to perform an overriding MERGE on the original object
// with your changes. Think of this as updating only fields you have defined and are non-zero.
// If you wish to specify a different strategy, the following ones are valid:
//	- "overwrite" will replace the entirety of original "foo" with your definition, discarding any previous state.
//	- "inherit" will apply a merge in reverse - merging the original "foo" into your definition, overwriting any fields.
//	- "panic" will raise a runtime error and prevent further execution. This can be a very helpful way to avoid state on "root" definitions.
type OnConflict struct {
	Do string `hcl:"do,attr" json:"do,omitempty"`
}

// CreateIndex maps out all of the known networks, identities, scripts, commands, and hosts within a laforge configuration snapshot.
func (l *Laforge) CreateIndex() {
	l.HostMap = map[string]*Host{}
	l.NetworkMap = map[string]*Network{}
	l.IdentityMap = map[string]*Identity{}
	l.ScriptMap = map[string]*Script{}
	l.CommandMap = map[string]*Command{}
	for _, x := range l.Hosts {
		l.HostMap[x.Hostname] = x
		x.Caller = l.Caller
	}
	for _, x := range l.Networks {
		l.NetworkMap[x.Name] = x
		x.Caller = l.Caller
	}
	for _, x := range l.Identities {
		l.IdentityMap[x.ID] = x
		x.Caller = l.Caller
	}
	for _, x := range l.Scripts {
		l.ScriptMap[x.Name] = x
		x.Caller = l.Caller
	}
	for _, x := range l.Commands {
		l.CommandMap[x.Name] = x
		x.Caller = l.Caller
	}
}

// Update performs a patching operation on source (l) with diff (diff), using the diff's merge conflict settings as appropriate.
func (l *Laforge) Update(diff *Laforge) (*Laforge, error) {
	l.Caller = diff.Caller
	if l.Filename != diff.Filename && diff.Filename != "" {
		l.Filename = diff.Filename
	}
	if l.BaseDir != diff.BaseDir && diff.BaseDir != "" {
		l.BaseDir = diff.BaseDir
	}
	var err error
	newUser := l.User
	err = mergo.Merge(&newUser, diff.User, mergo.WithOverride)
	if err != nil {
		return l, err
	}
	l.User = newUser
	newCompetition := l.Competition
	err = mergo.Merge(&newCompetition, diff.Competition, mergo.WithOverride)
	if err != nil {
		return l, err
	}
	l.Competition = newCompetition
	newEnvironment := l.Environment
	err = mergo.Merge(&newEnvironment, diff.Environment, mergo.WithOverride)
	if err != nil {
		return l, err
	}
	l.Environment = newEnvironment
	return l, err
}

// Mask attempts to apply a differential update betweeen base and layer, returning a modified base and any errors it encountered.
func Mask(base, layer *Laforge) (*Laforge, error) {
	layer.CreateIndex()
	var err error
	for name, obj := range layer.HostMap {
		orig, found := base.HostMap[name]
		if !found {
			base.HostMap[name] = obj
			continue
		}
		err = orig.Update(obj)
		if err != nil {
			return nil, err
		}
	}
	for name, obj := range layer.NetworkMap {
		orig, found := base.NetworkMap[name]
		if !found {
			base.NetworkMap[name] = obj
			continue
		}
		err = orig.Update(obj)
		if err != nil {
			return nil, err
		}
	}
	for name, obj := range layer.IdentityMap {
		orig, found := base.IdentityMap[name]
		if !found {
			base.IdentityMap[name] = obj
			continue
		}
		err = orig.Update(obj)
		if err != nil {
			return nil, err
		}
	}
	for name, obj := range layer.ScriptMap {
		orig, found := base.ScriptMap[name]
		if !found {
			base.ScriptMap[name] = obj
			continue
		}
		err = orig.Update(obj)
		if err != nil {
			return nil, err
		}
	}
	for name, obj := range layer.CommandMap {
		orig, found := base.CommandMap[name]
		if !found {
			base.CommandMap[name] = obj
			continue
		}
		err = orig.Update(obj)
		if err != nil {
			return nil, err
		}
	}
	return base.Update(layer)
}
