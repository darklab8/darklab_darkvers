/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autogit/actions"
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "next semantic version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", actions.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	actions.VersionDisableNewLine = versionCmd.PersistentFlags().Bool("no-newline", false, "Disable newline")

	actions.VersionDisableVFlag = versionCmd.PersistentFlags().Bool("no-v", false, "Disable v flag")
	actions.VersionBuildMeta = versionCmd.PersistentFlags().String("build", "", "Build metadata, not affecting semantic versioning. Added as semver+build")
	actions.VersionAlpha = versionCmd.PersistentFlags().Bool("alpha", false, "Enable next version as alpha")
	actions.VersionBeta = versionCmd.PersistentFlags().Bool("beta", false, "Enable next version as beta")
	actions.VersionPrerelease = versionCmd.PersistentFlags().Bool("rc", false, "Enable next version as prerelease")
}
