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

var logong2Backup = []string{
	boldb("                                               _______----_______               "),
	boldb("                                    ___---~~~~~.. ... .... ... ..~~~~~---___    "),
	boldb("                              _ ==============================================  "),
	boldb("  __________________________ - .. ..   _--~~~~~-------____-------~~~~~          "),
	boldb(" (______________________][__)____     -                                         "),
	boldb("    /       /______---~~~.. .. ..~~-_~                      [VERSION] 0.0.1     "),
	boldb("   <_______________________________-                         [AUTHOR] gen0cide  "),
	boldw(" dP    ~~~~~~~-----__     ^    __-                                              "),
	boldw(" 88                     _' `_      88888888b                                    "),
	boldw(" 88        .d8888b.  .-~'   `~-.   88      .d8888b. 88d888b. .d8888b. .d8888b.  "),
	boldw(" 88        88'  `88 (  ' __. `  ) a88aaaaa 88'  `88 88'  `88 88'  `88 88ooood8  "),
	boldw(" 88        88.  .88  `-''___``-'   88      88.  .88 88       88.  .88 88.  ...  "),
	boldw(" 88888888P `88888P8a   `     '     dP      `88888P' dP       `8888P88 `88888P'  "),
	normb(" ooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo~~~~.88~ooooooooo "),
	boldg(" github.com/gen0cide/laforge                                   d8888P           "),
	nocol("                                                                                "),
}

var logong2 = []string{
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
	nocol("                                                                              "),
	"",
}

var logong = []string{
	"",
	boldblue("                                               _______----_______                  "),
	boldblue("                                    ___---~~~~~.. ... .... ... ..~~~~~---___       "),
	boldblue("                              _ ==============================================     "),
	boldblue("  __________________________ - .. ..   _--~~~~~-------____-------~~~~~             "),
	boldblue(" (______________________][__)____     -                                            "),
	boldblue("    /       /______---~~~.. .. ..~~-_~                      %s%s%s %s        ", boldgreen("["), color.HiWhiteString("VERSION"), boldgreen("]"), boldred(Version)),
	boldblue("   <_______________________________-                         %s%s%s %s     ", boldgreen("["), color.HiWhiteString("AUTHOR"), boldgreen("]"), boldred(AuthorHandle)),
	boldwhite(" dP    %s                                                ", boldblue("~~~~~~~-----__           __-")),
	boldwhite(" 88                  %s %s                                         ", boldblue("~~~~~~~~~~~"), boldwhite("88888888b")),
	boldwhite(" 88        .d8888b.              88        .d8888b. 88d888b. .d8888b. .d8888b.     "),
	boldwhite(" 88        88'  `88             a88aaaaa   88'  `88 88'  `88 88'  `88 88ooood8     "),
	boldwhite(" 88        88.  .88              88        88.  .88 88       88.  .88 88.  ...     "),
	boldwhite(" 88888888P `88888P8a             dP        `88888P' dP       `8888P88 `88888P'     "),
	color.HiBlueString(" ooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo~~~~.%s%s    ", boldwhite("88"), color.HiBlueString("~ooooooooo")),
	fmt.Sprintf(" %s                                   %s              ", boldgreen(RepoURL), boldwhite("d8888P")),
	"                                                                                      ",
	"",
}

var logo = `
          ((((            ,(((/         ((((((((       ./(((/        ,(((((((,           *((((,        ((((((((
         @@@@          @@@@@@@@@       @@@@@@@@      @@@@@@@@@      ,@@@@@@@@@/       @@@@@@@@@,      @@@@@@@@
        @@@@          @@@/ .@@@      .@@@          .@@@  (@@@      &@@@  ,@@@(       @@@/  @@@*     .@@@
       @@@%%          @@@  .@@@      (@@@          .@@@  (@@@      &@@@  ,@@@.       @@@*  &.*      (@@@
      @@@%%         .@@@ *@@@@      (@@@ @@@@     (@@@  @@@@      &@@@ %%@@@@.       @@@*.@@@@      (@@@ @@@@
     @@@(         *@@@@@@@@@      %%@@@@@@@*     %%@@@  @@@@      @@@@@@@@#         @@@ (@@@@      %%@@@@@@@,
    @@@.         *@@@  %%@@@      @@@@/         %%@@@  @@@@      @@@@@@@@          @@@ ,*@@@      @@@@,
   @@@.         %%@@@  %%@@@      @@@@          @@@@  @@@&      @@@@ @@@@        *@@@  %%@@@      @@@@      %s%s
`
var logo2 = `  @@@.         %%@@@  @@@@      @@@@          @@@@  @@@&      @@@(  @@@@       *@@@  %%@@@      @@@@
,@@@@@@@@     @@@@  @@@@      @@@%%          @@@@@@@@@*      @@@*   @@@@      ,@@@@@@@@%%      @@@@@@@@*  /by/
                                                                                                        %s
                %s   %s   %s
                                     %s

`
var subTitle = `S E C U R I T Y    C O M P E T I T I O N    A U T O M A T I O N`
var url = `github.com/gen0cide/laforge`

// PrintLogo prints the Laforge logo to Stdout with color if possible.
func PrintLogo() {
	logoText := strings.Join(logong2, "\n")
	// logoText := fmt.Sprintf(
	// 	"%s%s",
	// 	color.HiBlueString(logo, boldyellow("v"), boldwhite(Version)),
	// 	color.HiBlueString(logo2, boldred(AuthorHandle), boldgreen(`--`), boldwhite(subTitle), boldgreen(`--`), boldyellow(url)),
	// )
	fmt.Fprint(color.Output, logoText)
}
