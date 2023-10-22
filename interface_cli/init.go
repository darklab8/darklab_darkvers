/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package interface_cli

import (
	"autogit/interface_cli/actions"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init repository settings. " + actions.InitAdvice,
	Run: func(cmd *cobra.Command, args []string) {
		shared.init.Run()
		actions.Initialize(*initGLobally)
	},
}

var initGLobally *bool

func init() {
	rootCmd.AddCommand(initCmd)
	initGLobally = initCmd.Flags().BoolP("global", "g", false, "Init settings file globally")
}