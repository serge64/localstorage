package cache_test

import (
	"testing"

	"github.com/serge64/localstorage/cache"
)

func TestCache(t *testing.T) {
	cache := cache.New()
	False(t, cache.Cached())
	cache.Save()
	True(t, cache.Cached())
	cache.Reset()
	False(t, cache.Cached())
}

func True(t *testing.T, b bool) {
	if !b {
		t.Error("Value expected 'true' but got 'false'.")
	}
}

func False(t *testing.T, b bool) {
	if b {
		t.Error("Value expected 'false' but got 'true'.")
	}
}
