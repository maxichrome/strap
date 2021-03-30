package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Failure logs error output and a description of when something went wrong
func Failure(description string, err error) {
	var (
		errLines = strings.Split(err.Error(), "\n")
	)

	fmt.Printf("Funky failure %s:\n", description)

	if len(errLines) > 1 {
		var numberColWidth int = len(fmt.Sprint((len(errLines) + 1)))
		for i, line := range errLines {
			fmt.Printf("%s %s\n", color.New(color.Faint).Add(color.BgBlack).Sprintf(" %*d ", numberColWidth, i + 1), line)
		}
	} else {
		fmt.Printf("  %s\n", errLines[0])
	}

	os.Exit(1)
}
