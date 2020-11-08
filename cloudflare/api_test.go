package cloudflare

import (
	"os"
	"testing"
)

func TestPurgeURLs(t *testing.T) {
	api := API{
		os.Getenv("TEST_CLOUDFLARE_TOKEN"),
		os.Getenv("TEST_CLOUDFLARE_ZONE_ID"),
	}
	if err := api.PurgeCache("https://cdn.my-shell.com/images/static/about/1988.jpg"); err != nil {
		t.Error(err)
	}
}
