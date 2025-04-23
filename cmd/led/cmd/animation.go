/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// animationCmd represents the animation command
var animationCmd = &cobra.Command{
	Use:   "animation [name]",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		ShowAnimation(args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	showCmd.AddCommand(animationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// animationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// animationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
