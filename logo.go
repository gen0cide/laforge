package laforge

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	boldgreen  = color.New(color.FgHiGreen, color.Bold).SprintFunc()
	boldwhite  = color.New(color.FgHiWhite, color.Bold).SprintFunc()
	boldred    = color.New(color.FgHiRed, color.Bold).SprintFunc()
	boldyellow = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	boldcyan   = color.New(color.FgHiCyan, color.Bold).SprintFunc()
	boldblue   = color.New(color.FgHiBlue, color.Bold).SprintFunc()
	plainblue  = color.New(color.FgHiBlue)
)

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
	logoText := fmt.Sprintf(
		"%s%s",
		color.HiBlueString(logo, boldyellow("v"), boldwhite(Version)),
		color.HiBlueString(logo2, boldred(AuthorHandle), boldgreen(`--`), boldwhite(subTitle), boldgreen(`--`), boldyellow(url)),
	)
	fmt.Fprint(color.Output, logoText)
}
