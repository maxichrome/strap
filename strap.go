package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var bar pb.ProgressBar

func main() {
	// todo: add version number here
	fmt.Println("strap  v{version}")

	bar = *pb.Simple.New(100)

	if _, err := os.Stat("~/.strap"); err == nil {
		// branch -> reinstall / check for changed files
	} else if os.IsNotExist(err) {
		firstRun()
	} else {
		// other funky failure
		fmt.Println("Unable to check for ~/.strap file - aborting!")
		return
	}
}

func firstRun() {
	fmt.Println("Let's get you set up! I'll need a bit of info from you...\n ")
	
	// todo: check for ~/.strap before re-prompting
	// (will have record of repo uri and stuff)
	var dotfileUri string

	for len(dotfileUri) < 1 {
		output, err := prompt("Enter your config repo uri: ")
		dotfileUri = output

		if err != nil {
			fmt.Println("Could not get config repo uri!", err)
			return
		}
	}

	bar.Start()

	for bar.Current() < 100 {
		bar.Increment()
		time.Sleep(time.Second / 2)
	}
	// todo: validate uri with git and proceed
}

func prompt(promptText string) (string, error) {
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
