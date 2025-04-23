//go:generate go run github.com/dmarkham/enumer -type=Dashboard -transform=kebab

package dashboard

type Dashboard int

const (
	Clock Dashboard = iota
	Autodarts
	Shopify
)
