package core

import (
	"fmt"
)

var (
	requiresTF = map[LFType]bool{
		LFTypeBuild:              true,
		LFTypeEnvironment:        true,
		LFTypeProvisionedHost:    true,
		LFTypeProvisionedNetwork: true,
		LFTypeTeam:               true,
	}

	requiresTFTaint = map[LFType]bool{
		LFTypeProvisionedHost: true,
	}
)

// CalculateTerraformNeeds attempts to determine what terraform steps must happen first
// for a given plan.
func CalculateTerraformNeeds(plan *Plan) (map[string][]string, error) {
	ret := map[string][]string{}
	teamsRequiringTFApply := map[string]bool{}
	for _, x := range plan.GlobalOrder {
		kind := TypeByPath(x)
		if _, ok := requiresTF[kind]; !ok {
			continue
		}
		teamID, err := GetTeamIDFromPath(x)
		if err != nil {
			continue
		}
		if ret[teamID] == nil || len(ret[teamID]) < 1 {
			ret[teamID] = []string{}
			ret[teamID] = append(ret[teamID], "init")
			ret[teamID] = append(ret[teamID], "refresh")
		}
		if kind == LFTypeProvisionedHost {
			host := plan.Graph.Metastore[x].Dependency
			phost, ok := host.(*ProvisionedHost)
			if !ok {
				continue
			}
			if phost.Conn == nil {
				ret[teamID] = append(ret[teamID], fmt.Sprintf("taint %s", phost.ID))
				teamsRequiringTFApply[teamID] = true
				continue // the host hasn't been deployed yet
			}
			ret[teamID] = append(ret[teamID], fmt.Sprintf("taint %s", phost.Conn.ResourceName))
			teamsRequiringTFApply[teamID] = true
		}
	}

	// now to clean up
	for tid := range teamsRequiringTFApply {
		ret[tid] = append(ret[tid], "apply -auto-approve -parallelism=10")
	}
	return ret, nil
}
