package main

import (
	"github.com/sanjid133/rest-user-store/version"
	"github.com/spf13/cobra"
)

// rootCmd is the root of all sub commands in the binary
// it doesn't have a Run method as it executes other sub commands
var rootCmd = &cobra.Command{
	Use:     "user",
	Short:   "user manages user",
	Version: version.Version,
}

// for now only requirement is to run server
func main() {
	rootCmd.Execute()
}

var cfgPath string

func init() {
	// register sub-commands here
	rootCmd.AddCommand(serveCmd)
}
