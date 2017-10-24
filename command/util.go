package command

import (
	"github.com/gen0cide/laforge/competition"
)

func InitConfig() (*competition.Competition, *competition.Environment) {
	competition.ValidateEnv()
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME Config: " + err.Error())
	}
	env := comp.CurrentEnv()
	if env == nil {
		competition.LogFatal("Cannot load environment! (Check ~/.lf_env)")
	}
	env.Competition = *comp
	return comp, env
}
