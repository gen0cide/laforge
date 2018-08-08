package laforge

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	boldgreen  = color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	boldwhite  = color.New(color.FgHiWhite, color.Bold).SprintfFunc()
	boldred    = color.New(color.FgHiRed, color.Bold).SprintfFunc()
	boldyellow = color.New(color.FgHiYellow, color.Bold).SprintfFunc()
	boldcyan   = color.New(color.FgHiCyan, color.Bold).SprintfFunc()
	boldb      = color.New(color.FgHiBlue, color.Bold).SprintfFunc()
	boldg      = color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	boldw      = color.New(color.FgHiWhite, color.Bold).SprintfFunc()
	boldr      = color.New(color.FgHiRed, color.Bold).SprintfFunc()
	boldy      = color.New(color.FgHiYellow, color.Bold).SprintfFunc()
	boldc      = color.New(color.FgHiCyan, color.Bold).SprintfFunc()
	boldm      = color.New(color.FgHiMagenta, color.Bold).SprintfFunc()
	britw      = color.New(color.FgHiWhite).SprintfFunc()
	normb      = color.New(color.FgBlue).SprintfFunc()
	nocol      = color.New(color.Reset).SprintfFunc()
	boldblue   = color.New(color.FgHiBlue, color.Bold).SprintfFunc()
	plainblue  = color.New(color.FgHiBlue)
)

// MonoLogo is a monochromatic representation of the laforge ASCII banner.
var MonoLogo = []string{
	"                                              _______----_______               ",
	"                                   ___---~~~~~.. ... .... ... ..~~~~~---___    ",
	"                             _ ==============================================  ",
	" __________________________ - .. ..   _--~~~~~-------____-------~~~~~          ",
	"(______________________][__)____     -                                         ",
	fmt.Sprintf("   /       /______---~~~.. .. ..~~-_~                      [VERSION] %s     ", Version),
	fmt.Sprintf("  <_______________________________-                         [AUTHOR] %s  ", AuthorHandle),
	"dP    ~~~~~~~-----__     ^    __-                                              ",
	"88                     _' `_      88888888b                                    ",
	"88        .d8888b.  .-~'   `~-.   88      .d8888b. 88d888b. .d8888b. .d8888b.  ",
	"88        88'  `88 (  ' __. `  ) a88aaaaa 88'  `88 88'  `88 88'  `88 88ooood8  ",
	"88        88.  .88  `-''___``-'   88      88.  .88 88       88.  .88 88.  ...  ",
	"88888888P `88888P8a   `     '     dP      `88888P' dP       `8888P88 `88888P'  ",
	"ooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo~~~~.88~ooooooooo ",
	fmt.Sprintf("%s                                   d8888P           ", RepoURL),
	"------------------------------------------------------------------------------ ",
}

// ColorLogo is a line slice of a colorized version of the logo
var ColorLogo = []string{
	nocol("------------------------------------------------------------------------------"),
	boldb("                                              _______%s%s               ", boldm("----"), boldb("_______")),
	boldb("                                   ___---~~~~~%s%s", boldw(".. ... .... ... .."), boldb("~~~~~---___    ")),
	boldb("                             _ ==============================================  "),
	boldb(" %s %s %s   %s          ", boldc("__________________________"), boldb("-"), boldw(".. .."), boldb("_--~~~~~-------%s%s", boldm("____"), boldb("-------~~~~~"))),
	boldb("%s%s                                         ", boldc("(______________________]%s", boldc("[__)")), boldb("____     -")),
	boldb("   /       /______---~~~%s%s                      %s%s%s %s     ", boldw(".. .. .."), boldb("~~-_~"), boldg("["), britw("VERSION"), boldg("]"), boldr(Version)),
	boldb("  <_______________________________-                         %s%s%s %s  ", boldg("["), britw("AUTHOR"), boldg("]"), boldr(AuthorHandle)),
	boldw("dP    %s     %s    %s                                              ", boldb("~~~~~~~-----__"), boldy("^"), boldb("__-")),
	boldw("88                     %s%s", boldy("_' `_"), boldw("      88888888b                                    ")),
	boldw("88        .d8888b.  %s%s", boldy(".-~'   `~-."), boldw("   88      .d8888b. 88d888b. .d8888b. .d8888b.  ")),
	boldw("88        88'  `88 %s%s", boldy("(  ' __. `  )"), boldw(" a88aaaaa 88'  `88 88'  `88 88'  `88 88ooood8  ")),
	boldw("88        88.  .88  %s%s", boldy("`-''___``-'"), boldw("   88      88.  .88 88       88.  .88 88.  ...  ")),
	boldw("88888888P `88888P8a   %s%s", boldy("`     '"), boldw("     dP      `88888P' dP       `8888P88 `88888P'  ")),
	normb("ooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo~~~~%s%s ", boldw(".88"), normb("~ooooooooo")),
	boldg("github.com/gen0cide/laforge                                   %s           ", boldw("d8888P")),
	nocol("------------------------------------------------------------------------------"),
	"",
}

// PrintLogo prints the Laforge logo to Stdout with color if possible.
func PrintLogo() {
	logoText := strings.Join(ColorLogo, "\n")
	fmt.Fprint(color.Output, logoText)
}
