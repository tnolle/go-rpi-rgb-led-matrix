/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/animation"
)

var listAnimation bool

// animationCmd represents the animation command
var animationCmd = &cobra.Command{
	Use:       "animation [name]",
	Short:     "A brief description of your command",
	ValidArgs: animation.AnimationStrings(),
	Run: func(cmd *cobra.Command, args []string) {
		if listAnimation {
			for _, name := range animation.AnimationValues() {
				fmt.Println(name)
			}
			return
		}
		if len(args) != 1 {
			_ = cmd.Usage()
			return
		}
		ShowAnimation(args[0])
	},
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}

func init() {
	showCmd.AddCommand(animationCmd)

	animationCmd.Flags().BoolVar(&listAnimation, "list", false, "List all available animations")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// animationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// animationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
