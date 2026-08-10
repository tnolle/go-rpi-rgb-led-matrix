package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

type catalogResource struct {
	name     string
	title    string
	endpoint string
}

var catalogResources = []catalogResource{
	{name: "image", title: "Images", endpoint: "/images"},
	{name: "gif", title: "GIFs", endpoint: "/gifs"},
	{name: "dashboard", title: "Dashboards", endpoint: "/dashboards"},
	{name: "animation", title: "Animations", endpoint: "/animations"},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List content available on the server",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		return printCatalog(catalogResources, true)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	for _, resource := range catalogResources {
		resource := resource
		listCmd.AddCommand(&cobra.Command{
			Use:   resource.name,
			Short: "List available " + resource.title,
			Args:  cobra.NoArgs,
			RunE: func(_ *cobra.Command, _ []string) error {
				return printCatalog([]catalogResource{resource}, false)
			},
		})
	}
}

func printCatalog(resources []catalogResource, grouped bool) error {
	configuredHosts := hosts()
	multipleHosts := len(configuredHosts) > 1
	var errs []error

	for hostIndex, host := range configuredHosts {
		if multipleHosts {
			if hostIndex > 0 {
				fmt.Println()
			}
			fmt.Printf("%s\n", host)
		}

		for resourceIndex, resource := range resources {
			if grouped {
				if resourceIndex > 0 {
					fmt.Println()
				}
				fmt.Printf("%s:\n", resource.title)
			}

			items, err := fetchCatalog(host, resource.endpoint)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %s: %w", host, resource.name, err))
				continue
			}
			if len(items) == 0 && grouped {
				fmt.Println("  (none)")
				continue
			}
			for _, item := range items {
				if grouped || multipleHosts {
					fmt.Printf("  %s\n", item)
				} else {
					fmt.Println(item)
				}
			}
		}
	}

	return errors.Join(errs...)
}
