package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hadenpf/strap/log"
	"github.com/hadenpf/strap/repo"
)

const (
	strapVer string = "0.1.0"
)

var (
	homeDir  string
	strapDir string

	configRepoURI string
)

func main() {
	fmt.Printf("strap v%s\n", strapVer)

	if userHomeDir, err := os.UserHomeDir(); err == nil {
		homeDir = userHomeDir
		strapDir = path.Join(homeDir, ".strap")
	} else {
		log.Failure("getting user home directory", err)
		return
	}

	// if _, err := os.Stat(path.Join(strapDir, "straplock.yml")); err == nil {
	if _, err := os.Stat(strapDir); err == nil {
		returningRun()
	} else if os.IsNotExist(err) {
		firstRun()
	} else {
		log.Failure("checking for ~/.strap", err)
		return
	}
}

func firstRun() {
	//! NOTE: DO NOT MODIFY FILESYSTEM UNTIL THE USER HAS PROVIDED A REPO URI!

	fmt.Println("Let's get you set up! I'll need a bit of info from you...\n ")

	if output, err := prompt("Enter your config repo uri:"); err == nil && len(configRepoURI) < 1 {
		//! past this line, filesystem modifications are OK.
		configRepoURI = output
	} else {
		log.Failure("getting config repo uri", err)
		return
	}

	if _, err := os.Stat(strapDir); err == nil || !os.IsNotExist(err) {
		if input, err := prompt("WARNING! This action will overwrite the existing strap sync.\nContinue? (Y/n)"); err == nil {
			if input == "Y" {
				os.RemoveAll(strapDir)
			} else if strings.ToLower(input) == "n" {
				fmt.Println("Smart choice.")
				return
			} else {

				return
			}
		} else {
			log.Failure("getting overwrite confirmation", err)
			return
		}
	}

	if err := repo.Clone(configRepoURI, path.Join(homeDir, ".strap/repo")); err == nil {
		fmt.Println("Uh, I think that worked?!")
	} else {
		log.Failure("cloning config repo", err)
		return
	}
}

func returningRun() {
	if input, err := prompt("Continue using existing repo? (y/n/?)"); err == nil {
		if strings.ToLower(input) == "y" {
			fmt.Println("Checking for updates...")

			// todo: diff tracked dotfiles against repo
			//:
			// copy tracked dotfiles into ghost-repo path
			// use lockfile to track respective paths
			// 
			// then `git status` in ghost repo to find
			// any changes | parse status output and
			// print it nicely, asking user to update the
			// file in ghost repo:
			//   Add this file as change? ((y)es/(n)o/(d)iff/(a)ll yes/(m)anual -> cd into repo)
			//
			// === FILE REVIEW COMPLETED ===
			//  >Commit changes? (Y/n/m)
			//  >Push changes (Y/n/(F)orce/m)
		} else if strings.ToLower(input) == "n" {
			firstRun()
		} else if input == "?" {
			fmt.Println(configRepoURI)
			returningRun()
			return
		} else {
			fmt.Println("Invalid entry - try again")
			returningRun()
		}
	} else {
		log.Failure("getting continue input", err)
		return
	}
}

func prompt(promptText string) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	timesPrompted := 0

	for len(scanner.Text()) < 1 && timesPrompted < 5 {
		timesPrompted++

		fmt.Printf("%s ", promptText)
		scanner.Scan()
	}

	if timesPrompted >= 5 {
		return "", errors.New("retry limit exceeded")
	}

	return scanner.Text(), nil
}

// StrapLock is the format of the straplock.yml file.
type StrapLock struct {
	// version of Strap the file was created with
	strapVersion string

	repo struct {
		uri        string
		headCommit string
		ghostDir   string
	}
}
