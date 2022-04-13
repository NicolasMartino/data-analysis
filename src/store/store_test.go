package store_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/NicolasMartino/data-analysis/src/models"
	"github.com/NicolasMartino/data-analysis/src/store"
	"github.com/stretchr/testify/require"
)

func TestStoreReadWrite(t *testing.T) {
	cache := store.NewSafeMap(1 * time.Second)

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("test %o", i)
		value := fmt.Sprintf("test %o-1", i)
		cache.Store(key, models.CacheUrlInfo{UrlInfo: models.Data{Body: value}})

		retrievedValue, ok := cache.Load(key)
		require.True(t, ok)
		require.Equal(t, value, retrievedValue.UrlInfo.Body)
	}

}

func TestNeverReadFromCache(t *testing.T) {
	cache := store.NewSafeMap(0)

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("test %o", i)
		value := fmt.Sprintf("test %o-1", i)
		cache.Store(key, models.CacheUrlInfo{UrlInfo: models.Data{Body: value}})

		retrievedValue, ok := cache.Load(key)
		require.False(t, ok)
		require.Empty(t, retrievedValue.UrlInfo.RequestUrl)
		require.Empty(t, retrievedValue.UrlInfo.Status)
		require.Empty(t, retrievedValue.UrlInfo.Body)
	}
}
