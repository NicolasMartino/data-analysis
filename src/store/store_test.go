package store_test

import (
	"fmt"
	"testing"

	"github.com/NicolasMartino/data-analysis/src/models"
	"github.com/NicolasMartino/data-analysis/src/store"
	"github.com/stretchr/testify/require"
)

// A test that shows that simple maps are unsafe
func testUnsafeStore(t *testing.T) {
	cache := make(map[string]models.CacheUrlInfo)

	for i := 0; i < 5000; i++ {
		value := fmt.Sprintf("test %o", i)
		cache[fmt.Sprintf("%o", i)] = models.CacheUrlInfo{UrlInfo: models.UrlInfo{Body: value}}

		go func(cache map[string]models.CacheUrlInfo, i int) {
			value := fmt.Sprintf("test %o", i)
			cache[fmt.Sprintf("%o", i)] = models.CacheUrlInfo{UrlInfo: models.UrlInfo{Body: value}}
		}(cache, i)

		require.Equal(t, value, cache[fmt.Sprintf("%o", i)].UrlInfo.Body)
	}
}
func TestStoreConcurrentWrite(t *testing.T) {
	cache := store.NewSafeMap()

	for i := 0; i < 5000; i++ {
		value := fmt.Sprintf("test %o", i)
		cache.Store(fmt.Sprintf("%o", i), models.CacheUrlInfo{UrlInfo: models.UrlInfo{Body: value}})

		go func(cache *store.SafeMap, i int) {
			value := fmt.Sprintf("test %o", i)
			cache.Store(fmt.Sprintf("%o", i), models.CacheUrlInfo{UrlInfo: models.UrlInfo{Body: value}})
		}(cache, i)

		retrievedValue, err := cache.Load(fmt.Sprintf("%o", i))
		require.Equal(t, true, err)
		require.Equal(t, value, retrievedValue.UrlInfo.Body)
	}

}
