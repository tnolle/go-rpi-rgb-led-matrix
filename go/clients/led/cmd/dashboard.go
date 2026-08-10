package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/dashboard"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:       "dashboard [name]",
	Short:     "Set the displayed dashboard",
	ValidArgs: dashboard.DashboardStrings(),
	RunE: func(_ *cobra.Command, args []string) error {
		return SetDashboard(args[0])
	},
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}

func init() {
	setCmd.AddCommand(dashboardCmd)
}
