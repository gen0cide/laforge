package core

// IdentitiesWithVarEquals allows you to retrieve all identites who var key matches value
func (l *Laforge) IdentitiesWithVarEquals(key, value string) []*Identity {
	ret := []*Identity{}
	for _, x := range l.Identities {
		if x.Vars[key] == value {
			ret = append(ret, x)
		}
	}

	return ret
}

// UniqIdentityVarValues allows you to gather unique var values from within identities
func (l *Laforge) UniqIdentityVarValues(key string) []string {
	ret := map[string]bool{}
	for _, x := range l.Identities {
		if val, ok := x.Vars[key]; ok {
			ret[val] = true
		}
	}
	final := []string{}
	for k := range ret {
		final = append(final, k)
	}

	return final
}
