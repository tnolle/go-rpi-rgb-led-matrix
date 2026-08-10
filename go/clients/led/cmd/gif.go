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
	Use:   "gif [name]",
	Short: "Set the displayed GIF",
	RunE: func(_ *cobra.Command, args []string) error {
		return SetGIF(args[0], once)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	setCmd.AddCommand(gifCmd)
	gifCmd.Flags().BoolVarP(&once, "once", "o", false, "Display the GIF only once")
}
