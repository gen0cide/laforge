package core

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// StatusMap gives a visual representation of the current state context
func StatusMap(curr StateContext) string {
	var l0, l1, l2, l3, l4 string
	if curr == TeamContext {
		l0 = fmt.Sprintf(" (%s) %s %s", boldc("%d", int(TeamContext)), boldw("*CURRENT*"), boldc(TeamContext.String()))
	} else {
		l0 = fmt.Sprintf(" (%s) %s", color.HiCyanString("%d", int(TeamContext)), color.HiCyanString(TeamContext.String()))
	}
	if curr == BuildContext {
		l1 = fmt.Sprintf(" (%s) %s %s", boldg("%d", int(BuildContext)), boldw("*CURRENT*"), boldg(BuildContext.String()))
	} else {
		l1 = fmt.Sprintf(" (%s) %s", color.HiGreenString("%d", int(BuildContext)), color.HiGreenString(BuildContext.String()))
	}
	if curr == EnvContext {
		l2 = fmt.Sprintf(" (%s) %s %s", boldy("%d", int(EnvContext)), boldw("*CURRENT*"), boldy(EnvContext.String()))
	} else {
		l2 = fmt.Sprintf(" (%s) %s", color.HiYellowString("%d", int(EnvContext)), color.HiYellowString(EnvContext.String()))
	}
	if curr == BaseContext {
		l3 = fmt.Sprintf(" (%s) %s %s", boldr("%d", int(BaseContext)), boldw("*CURRENT*"), boldr(BaseContext.String()))
	} else {
		l3 = fmt.Sprintf(" (%s) %s", color.HiRedString("%d", int(BaseContext)), color.HiRedString(BaseContext.String()))
	}
	if curr == GlobalContext {
		l4 = fmt.Sprintf(" (%s) %s %s", boldb("%d", int(GlobalContext)), boldw("*CURRENT*"), boldb(GlobalContext.String()))
	} else {
		l4 = fmt.Sprintf(" (%s) %s", color.HiBlueString("%d", int(GlobalContext)), color.HiBlueString(GlobalContext.String()))
	}
	return strings.Join([]string{l0, l1, l2, l3, l4}, "\n")
}
