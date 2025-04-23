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
		selected := viper.GetInt("selectedHost")
		if len(hosts) == 0 {
			fmt.Println("No hosts configured.")
		} else {
			fmt.Println("Configured hosts:")
			if selected == 0 {
				fmt.Println("0: All hosts *")
			} else {
				fmt.Println("0: All hosts")
			}
			for i, h := range hosts {
				if selected == i+1 {
					fmt.Printf("%d: %s *\n", i+1, h)
				} else {
					fmt.Printf("%d: %s\n", i+1, h)
				}
			}
		}
	},
	Args: cobra.ExactArgs(1),
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

var hostsDeleteCmd = &cobra.Command{
	Use:   "delete [index]",
	Short: "Delete a host by index",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var index int
		_, err := fmt.Sscanf(args[0], "%d", &index)
		if err != nil || index <= 0 {
			fmt.Println("Invalid index. Use `hosts` to view the list.")
			os.Exit(1)
		}
		hosts := viper.GetStringSlice("hosts")
		if index > len(hosts) {
			fmt.Println("Index out of range.")
			os.Exit(1)
		}
		hosts = append(hosts[:index-1], hosts[index:]...)
		viper.Set("hosts", hosts)
		if viper.GetInt("selectedHost") == index {
			viper.Set("selectedHost", 0) // Reset to "all hosts" if the selected one was deleted
		}
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error writing config:", err)
			os.Exit(1)
		}
		fmt.Println("Host deleted.")
	},
}

var hostsSelectCmd = &cobra.Command{
	Use:   "select [index]",
	Short: "Select a host by index (0 = all hosts)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var index int
		_, err := fmt.Sscanf(args[0], "%d", &index)
		if err != nil || index < 0 {
			fmt.Println("Invalid index.")
			os.Exit(1)
		}
		hosts := viper.GetStringSlice("hosts")
		if index > len(hosts) {
			fmt.Println("Index out of range.")
			os.Exit(1)
		}
		viper.Set("selectedHost", index)
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error writing config:", err)
			os.Exit(1)
		}
		fmt.Println("Host selected.")
	},
}

func init() {
	rootCmd.AddCommand(hostsCmd)
	hostsCmd.AddCommand(hostsAddCmd)
	hostsCmd.AddCommand(hostsDeleteCmd)
	hostsCmd.AddCommand(hostsSelectCmd)

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
