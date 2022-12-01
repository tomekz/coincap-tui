package cmd

import (
	"github.com/spf13/cobra"
)

type exportCmd struct {
	cmd *cobra.Command
}

func newExportCmd() *exportCmd {
	cmd := &cobra.Command{
		Use:     "report",
		Aliases: []string{"r"},
		Short:   "Print a markdown report of the given project to STDOUT",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			// TODO implement
			// ## Export
			//
			// You can extract a markdown file by running:
			//
			// ```sh
			// tt report
			// ```
			// It will output the given project (via `-p PROJECT`) to `STDOUT`. You can
			// then save it to a file, pipe to another software or do whatever you like:

			return nil
		},
	}

	return &exportCmd{cmd: cmd}
}
