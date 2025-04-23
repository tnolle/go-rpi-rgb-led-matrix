package cmd

import (
	"github.com/spf13/cobra"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use: "dashboard [name]",
	Run: func(cmd *cobra.Command, args []string) {
		ShowDashboard(args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	showCmd.AddCommand(dashboardCmd)
}
