package repo

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// Clone is a function which clones a git repo
func Clone(repoURI string, destination string) (error) {
	cloneCmd := exec.Command("git", "clone", repoURI, destination)
	
	stdOut := bufio.NewWriter(os.Stdout)
	errBuf := bytes.NewBuffer([]byte("a"))

	// cloneCmd.Stdin = 
	cloneCmd.Stdout = stdOut
	cloneCmd.Stderr = errBuf

	if err := cloneCmd.Run(); err != nil || len(errBuf.Bytes()) > 0 {
		return fmt.Errorf("%s%s", errBuf.String(), err.Error())
	}

	return nil
}
