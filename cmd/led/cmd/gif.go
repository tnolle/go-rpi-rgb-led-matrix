/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var once bool

// gifCmd represents the gif command
var gifCmd = &cobra.Command{
	Use: "gif [name]",
	Run: func(cmd *cobra.Command, args []string) {
		ShowGIF(args[0], once)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	showCmd.AddCommand(gifCmd)
	gifCmd.Flags().BoolVarP(&once, "once", "o", false, "Show the GIF only once")
}
