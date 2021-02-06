package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var bar pb.ProgressBar
var homeDir string
var configRepoUri string

func main() {
	// todo: add version number here
	fmt.Println("strap  v{version}")

	barTemplate := `{{ bar . "[" "#" "#" "~" "]"}} {{percent .}}`
	bar = *pb.ProgressBarTemplate(barTemplate).New(100)

	if userHomeDir, err := os.UserHomeDir(); err == nil {
		homeDir = userHomeDir
	} else {
		fmt.Println("Funky failure getting user home directory - aborting!", err)
		return
	}

	if _, err := os.Stat(path.Join(homeDir, ".strap.yml")); err == nil {
		returningRun()
	} else if os.IsNotExist(err) {
		firstRun()
	} else {
		fmt.Println("Funky failure checking for ~/.strap file - aborting!")
		return
	}
}

func firstRun() {
	fmt.Println("Let's get you set up! I'll need a bit of info from you...\n ")
	
	// todo: check ~/.strap contents before re-prompting for uri
	// (will have record of repo uri and stuff)

	if output, err := Prompt("Enter your config repo uri: "); err == nil && len(configRepoUri) < 1 {
		configRepoUri = output
	} else {
		fmt.Println(err)
		fmt.Println("Funky failure getting config repo uri - aborting!")
		return
	}

	bar.Start()

	for bar.Current() < 100 {
		bar.Increment()
		time.Sleep(time.Second / 10)
	}

	bar.Finish()
	// todo: validate uri with git and proceed
}

func returningRun() {
	 if input, err := Prompt("Continue using existing repo? (y/n/?) "); err == nil {
		if strings.ToLower(input) == "y" {
			fmt.Println("Checking for updates...")

			// todo: diff tracked dotfiles against repo
			// ignore apps at this point
		} else if strings.ToLower(input) == "n" {
			firstRun()
		} else if input == "?" {
			fmt.Println(configRepoUri)
			returningRun()
			return
		} else {
			fmt.Println("Invalid entry - try again")
			returningRun()
		}
	} else {
		fmt.Println(err)
		fmt.Println("Funky failure getting continue input - aborting!")
		return
	}
}

func Prompt(promptText string) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var timesPrompted int = 0

	for len(scanner.Text()) < 1 && timesPrompted < 5 {
		timesPrompted++

		fmt.Print(promptText)
		scanner.Scan()
	}

	if timesPrompted >= 5 {
		return "", errors.New("Retry limit exceeded")
	}

	return scanner.Text(), nil
}

type StrapLock struct {
	strap   string // version

	repo struct {
		uri					string
		head_commit	string
    ghost_dir   string
	}
}
