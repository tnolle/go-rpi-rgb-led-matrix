//go:generate go run github.com/dmarkham/enumer@v1.6.3 -type=Dashboard -transform=kebab

package dashboard

type Dashboard int

const (
	Clock Dashboard = iota
	Autodarts
	Shopify
)
