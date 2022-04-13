package store_test

import (
	"fmt"
	"testing"

	"github.com/NicolasMartino/data-analysis/src/models"
	"github.com/NicolasMartino/data-analysis/src/store"
	"github.com/stretchr/testify/require"
)

func TestStoreReadWrite(t *testing.T) {
	cache := store.NewSafeMap()

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("test %o", i)
		value := fmt.Sprintf("test %o-1", i)
		cache.Store(key, models.CacheUrlInfo{UrlInfo: models.Data{Body: value}})

		retrievedValue, ok := cache.Load(key)
		require.True(t, ok)
		require.Equal(t, value, retrievedValue.UrlInfo.Body)
	}

}
