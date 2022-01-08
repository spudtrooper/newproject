package gitversion

import (
	"flag"
	"fmt"
)

var (
	version = flag.Bool("version", false, "print version")
)

// Prints the big version if --version is set and returns true, otherwise reutrns false
func CheckVersionFlag() bool {
	if *version {
		fmt.Printf("Version: %s\n", theGitVersion)
		return true
	}
	return false
}
