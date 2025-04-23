/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use: "hosts",
	Run: func(cmd *cobra.Command, args []string) {
		hosts := viper.GetStringSlice("hosts")
		if len(hosts) == 0 {
			fmt.Println("No hosts configured.")
			return
		}
		fmt.Println("Configured hosts:")
		for _, h := range hosts {
			fmt.Println("-", h)
		}
	},
}

var hostsAddCmd = &cobra.Command{
	Use:   "add [host]",
	Short: "Add a host to the config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		hosts := viper.GetStringSlice("hosts")
		hosts = append(hosts, host)
		viper.Set("hosts", hosts)
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error writing config:", err)
			os.Exit(1)
		}
		fmt.Println("Host added:", host)
	},
}

func init() {
	rootCmd.AddCommand(hostsCmd)
	hostsCmd.AddCommand(hostsAddCmd)

	viper.SetConfigName("led")  // name of config file (without extension)
	viper.SetConfigType("toml") // config file type
	cfgDir, _ := os.UserConfigDir()
	viper.AddConfigPath(cfgDir) // look for config in the current directory

	// If the config file doesn't exist, create it
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create it
			if err := viper.WriteConfigAs(cfgDir + "/led.toml"); err != nil {
				fmt.Println("Error creating config file:", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Error reading config:", err)
			os.Exit(1)
		}
	}
}
