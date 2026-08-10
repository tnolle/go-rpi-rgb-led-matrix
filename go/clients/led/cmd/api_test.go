package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestFetchCatalog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/images" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"items":["alpha","zebra"]}`)
	}))
	defer server.Close()

	items, err := fetchCatalog(server.URL, "/images")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"alpha", "zebra"}
	if !reflect.DeepEqual(items, want) {
		t.Fatalf("unexpected catalog: got %v, want %v", items, want)
	}
}

func TestDoEscapesName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("name"); got != "name with spaces & symbols" {
			t.Fatalf("unexpected name: %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	if err := do(server.URL+"/image", "name with spaces & symbols"); err != nil {
		t.Fatal(err)
	}
}

func TestHostsDefaultsToLocalServer(t *testing.T) {
	previousHosts := viper.GetStringSlice("hosts")
	previousSelected := viper.GetInt("selectedHost")
	t.Cleanup(func() {
		viper.Set("hosts", previousHosts)
		viper.Set("selectedHost", previousSelected)
	})

	viper.Set("hosts", []string{})
	viper.Set("selectedHost", 0)
	if got := hosts(); !reflect.DeepEqual(got, []string{HOST}) {
		t.Fatalf("unexpected hosts: %v", got)
	}
}
