package version

import (
	"fmt"
	"strings"
)

var (
	// The git commit that was compiled. These will be filled in by the
	// compiler.
	GitCommit   string
	GitDescribe string

	// The main version number that is being run at the moment.
	//
	// Version must conform to the format expected by github.com/hashicorp/go-version
	// for tests to work.
	Version = "1.2.0"

	// A pre-release marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = "beta3"
)

// GetHumanVersion composes the parts of the version in a way that's suitable
// for displaying to humans.
func GetHumanVersion() string {
	version := Version
	if GitDescribe != "" {
		version = GitDescribe
	}

	release := VersionPrerelease
	if GitDescribe == "" && release == "" {
		release = "dev"
	}
	if release != "" {
		version += fmt.Sprintf("-%s", release)
		if GitCommit != "" {
			version += fmt.Sprintf(" (%s)", GitCommit)
		}
	}

	// Strip off any single quotes added by the git information.
	return strings.Replace(version, "'", "", -1)
}
