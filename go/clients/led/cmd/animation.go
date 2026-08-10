/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/animation"
)

// animationCmd represents the animation command
var animationCmd = &cobra.Command{
	Use:       "animation [name]",
	Short:     "Set the displayed animation",
	ValidArgs: animation.AnimationStrings(),
	RunE: func(_ *cobra.Command, args []string) error {
		return SetAnimation(args[0])
	},
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}

func init() {
	setCmd.AddCommand(animationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// animationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// animationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
