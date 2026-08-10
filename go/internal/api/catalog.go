package api

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/animation"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/renderers/dashboard"
)

type catalog struct {
	imagesDir string
	gifsDir   string
}

func defaultCatalog() catalog {
	return catalog{
		imagesDir: "images/pngs",
		gifsDir:   "images/gifs",
	}
}

func (c catalog) images() ([]string, error) {
	return assetNames(c.imagesDir, ".png")
}

func (c catalog) gifs() ([]string, error) {
	return assetNames(c.gifsDir, ".gif")
}

func (c catalog) dashboards() []string {
	items := dashboard.DashboardStrings()
	sort.Strings(items)
	return items
}

func (c catalog) animations() []string {
	items := animation.AnimationStrings()
	sort.Strings(items)
	return items
}

func contains(items []string, name string) bool {
	index := sort.SearchStrings(items, name)
	return index < len(items) && items[index] == name
}

func assetNames(dir, extension string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}

	items := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if !strings.EqualFold(ext, extension) {
			continue
		}
		items = append(items, strings.TrimSuffix(entry.Name(), ext))
	}
	sort.Strings(items)
	return items, nil
}
