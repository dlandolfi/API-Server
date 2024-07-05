package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

// Constants for test data
const (
	TestAPIKey      = "test-api-key"
	TestRedisPW     = "test-password"
	TestPriceURL    = "http://example.com/price"
	TestNewsFeedURL = "http://example.com/news"
	TestRedisKey    = "priceObject"
	TestRedisValue  = "test-data"
)

func TestFetchURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	t.Run("fetchURL should return correct response", func(t *testing.T) {
		result, err := fetchURL(server.URL, map[string]string{"key": "value"})
		assert.NoError(t, err)
		assert.Contains(t, result, `"status": "ok"`)
	})
}

func TestStoreData(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	ctx := context.Background()

	t.Run("storeData should store data in Redis", func(t *testing.T) {
		mock.ExpectSet(TestRedisKey, TestRedisValue, 0).SetVal("OK")

		err := storeData(rdb, ctx, TestRedisKey, TestRedisValue)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLoadConfig(t *testing.T) {
	content := []byte(`{
		"metals_api_key": "` + TestAPIKey + `",
		"redis_pw": "` + TestRedisPW + `",
		"PriceURL": "` + TestPriceURL + `",
		"NewsFeedURL": "` + TestNewsFeedURL + `"
	}`)
	tmpfile, err := os.CreateTemp("", "config-test.json")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	assert.NoError(t, err)
	tmpfile.Close()

	t.Run("loadConfig should load configuration correctly", func(t *testing.T) {
		config, err := loadConfig(tmpfile.Name())
		assert.NoError(t, err)
		assert.Equal(t, TestAPIKey, config.APIKey)
		assert.Equal(t, TestPriceURL, config.PriceURL)
		assert.Equal(t, TestNewsFeedURL, config.NewsFeedURL)
		assert.Equal(t, TestRedisPW, config.REDISPW)
	})
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
