package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/dashboard"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:       "dashboard [name]",
	ValidArgs: dashboard.DashboardStrings(),
	Run: func(cmd *cobra.Command, args []string) {
		ShowDashboard(args[0])
	},
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}

func init() {
	showCmd.AddCommand(dashboardCmd)
}
