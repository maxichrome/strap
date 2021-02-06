package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

const (
	strapVer string = "0.1.0"
)

var (
	bar pb.ProgressBar

	homeDir string
	strapDir string

	configRepoUri string
)

func main() {
	// todo: add version number here
	fmt.Printf("strap  v%s\n", strapVer)

	barTemplate := `{{ bar . "[" "#" "#" "~" "]"}} {{percent .}}`
	bar = *pb.ProgressBarTemplate(barTemplate).New(100)

	if userHomeDir, err := os.UserHomeDir(); err == nil {
		homeDir = userHomeDir
		strapDir = path.Join(homeDir, ".strap")
	} else {
		fmt.Println("Funky failure getting user home directory - aborting!", err)
		return
	}

	// if _, err := os.Stat(path.Join(strapDir, "straplock.yml")); err == nil {
	if _, err := os.Stat(strapDir); err == nil {
		returningRun()
	} else if os.IsNotExist(err) {
		firstRun()
	} else {
		fmt.Println(err)
		fmt.Println("Funky failure checking for ~/.strap directory - aborting!")
		return
	}
}

func firstRun() {
	//! NOTE: DO NOT MODIFY FILESYSTEM UNTIL THE USER HAS PROVIDED A REPO URI!

	fmt.Println("Let's get you set up! I'll need a bit of info from you...\n ")

	//! past this line, filesystem modifications are OK.
	if output, err := Prompt("Enter your config repo uri:"); err == nil && len(configRepoUri) < 1 {
		configRepoUri = output
	} else {
		fmt.Println(err)
		fmt.Println("Funky failure getting config repo uri - aborting!")
		return
	}

	fmt.Printf("Cloning from %s...\n", configRepoUri)

	if _, err := os.Stat(strapDir); err == nil || !os.IsNotExist(err) {
		if input, err := Prompt("\n\nTHIS IS NOT RECOVERABLE:\nAre you sure you would like to remove your old strap synchronization? (Y/n)"); err == nil {
			if input == "Y" {
				os.RemoveAll(strapDir)
			} else {
				fmt.Println("You really scared me for a minute! Have a nice day...")
				return
			}
		}
	}

	cloneCmd := exec.Command("git", "clone", configRepoUri, path.Join(homeDir, ".strap/repo"))

	if out, err := cloneCmd.Output(); err == nil {
		fmt.Print(out)
		
		fmt.Println("Uh, I think that worked?!")
	} else {
		fmt.Println(err)
		fmt.Println("Funky failure while cloning config repo")
		return
	}
}

func returningRun() {
	 if input, err := Prompt("Continue using existing repo? (y/n/?)"); err == nil {
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

		fmt.Printf("%s ", promptText)
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
