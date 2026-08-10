package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/dashboard"
)

var listDashboard bool

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:       "dashboard [name]",
	ValidArgs: dashboard.DashboardStrings(),
	Run: func(cmd *cobra.Command, args []string) {
		if listDashboard {
			for _, name := range dashboard.DashboardValues() {
				fmt.Println(name)
			}
			return
		}
		if len(args) != 1 {
			_ = cmd.Usage()
			return
		}
		ShowDashboard(args[0])
	},
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}

func init() {
	showCmd.AddCommand(dashboardCmd)
	dashboardCmd.Flags().BoolVar(&listDashboard, "list", false, "List all available dashboards")
}
