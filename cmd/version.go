package cmd

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var version string

// Version returns coincap-tui version
func Version() string {
	ver := "(devel)"
	if version != "" {
		ver = version
	} else if buildInfo, ok := debug.ReadBuildInfo(); ok {
		ver = buildInfo.Main.Version
	}

	if !strings.HasPrefix(ver, "v") {
		ver = fmt.Sprintf("v%s", ver)
	}

	return ver
}

// VersionCmd ...
func VersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Displays the current version",
		Long:  `The version command displays the current version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version())
		},
	}

	return versionCmd
}
