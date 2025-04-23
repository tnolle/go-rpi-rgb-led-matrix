package rgbmatrix

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Auth struct {
		ClientID     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	} `toml:"auth"`
	Shopify struct {
		TotalSales   string `toml:"total_sales"`
		TotalOrders  string `toml:"total_orders"`
		MonthlySales string `toml:"monthly_sales"`
		TodaySales   string `toml:"today_sales"`
		TodayOrders  string `toml:"today_orders"`
	}
	Dashboards struct {
		Font string `toml:"font"`
	} `toml:"dashboards"`
	Options        MatrixOptions  `toml:"options"`
	RuntimeOptions RuntimeOptions `toml:"runtime_options"`
}

func LoadConfig() Config {
	path := "./config.toml"

	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("No config file found", err)
	}

	var config Config
	err = toml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal("Error parsing config file", err)
	}

	// Defaults
	if config.Options.PWMBits == 0 {
		config.Options.PWMBits = 11
	}
	if config.Options.PWMLSBNanoseconds == 0 {
		config.Options.PWMLSBNanoseconds = 130
	}
	if config.Options.PWMDitherBits == 0 {
		config.Options.PWMDitherBits = 0
	}
	if config.Options.ScanMode == 0 {
		config.Options.ScanMode = Progressive
	}
	if config.Options.ChainLength == 0 {
		config.Options.ChainLength = 1
	}
	if config.Options.Parallel == 0 {
		config.Options.Parallel = 1
	}
	if config.Dashboards.Font == "" {
		config.Dashboards.Font = "fonts/7x14.bdf"
	}

	return config
}
