package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string
)


var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  `Print the version of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print the version
    fmt.Println("Version:", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
