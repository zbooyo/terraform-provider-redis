package redisprovider

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"redis": testAccProvider,
	}
}

// ----------- Unit tests with redismock ------------

func TestDataSourceRedisStringRead_Success(t *testing.T) {
	// Create mock client
	client, mock := redismock.NewClientMock()
	cfg := &ProviderConfig{RedisClient: client}
	ctx := context.Background()

	// Setup data with a key and max_wait_seconds
	d := schema.TestResourceDataRaw(t, dataSourceRedisString().Schema, map[string]interface{}{
		"key":              "mock-key",
		"max_wait_seconds": 1,
	})

	// Mock Redis GET returning a value on first try
	mock.ExpectGet("mock-key").SetVal("mock-value")

	diags := dataSourceRedisStringRead(ctx, d, cfg)
	require.Len(t, diags, 0)
	assert.Equal(t, "mock-key", d.Id())
	assert.Equal(t, "mock-key", d.Get("key"))
	assert.Equal(t, "mock-value", d.Get("value"))

	// Ensure all expectations met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDataSourceRedisStringRead_WaitThenSuccess(t *testing.T) {
	client, mock := redismock.NewClientMock()
	cfg := &ProviderConfig{RedisClient: client}
	ctx := context.Background()

	d := schema.TestResourceDataRaw(t, dataSourceRedisString().Schema, map[string]interface{}{
		"key":              "delayed-key",
		"max_wait_seconds": 2,
	})

	// Simulate Redis GET returns empty or redis.Nil first 2 times, then value
	mock.ExpectGet("delayed-key").RedisNil()
	mock.ExpectGet("delayed-key").SetVal("")
	mock.ExpectGet("delayed-key").SetVal("final-value")

	diags := dataSourceRedisStringRead(ctx, d, cfg)
	require.Len(t, diags, 0)
	assert.Equal(t, "delayed-key", d.Id())
	assert.Equal(t, "final-value", d.Get("value"))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDataSourceRedisStringRead_Timeout(t *testing.T) {
	client, mock := redismock.NewClientMock()
	cfg := &ProviderConfig{RedisClient: client}
	ctx := context.Background()

	d := schema.TestResourceDataRaw(t, dataSourceRedisString().Schema, map[string]interface{}{
		"key":              "missing-key",
		"max_wait_seconds": 1,
	})

	// All Redis GETs return redis.Nil (key missing)
	// Expect at least 2 calls due to 500ms sleep in loop
	mock.ExpectGet("missing-key").RedisNil()
	mock.ExpectGet("missing-key").RedisNil()

	diags := dataSourceRedisStringRead(ctx, d, cfg)
	require.Len(t, diags, 0)
	assert.Equal(t, "", d.Id())
	assert.Equal(t, "", d.Get("value"))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDataSourceRedisStringRead_RedisError(t *testing.T) {
	client, mock := redismock.NewClientMock()
	cfg := &ProviderConfig{RedisClient: client}
	ctx := context.Background()

	d := schema.TestResourceDataRaw(t, dataSourceRedisString().Schema, map[string]interface{}{
		"key":              "error-key",
		"max_wait_seconds": 1,
	})

	mock.ExpectGet("error-key").SetErr(fmt.Errorf("redis connection error"))

	diags := dataSourceRedisStringRead(ctx, d, cfg)
	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, "redis connection error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ----------- Acceptance test with real Redis ------------

func TestAccDataSourceRedisString_Basic(t *testing.T) {
	testAccPreCheck(t)

	resourceName := "data.redis_string.test"
	key := fmt.Sprintf("tf_acc_ds_key_%d", time.Now().UnixNano())
	value := "tf_acc_ds_value"

	client := newRealRedisClient(t)
	defer client.Close()

	// Set key manually before test
	err := client.Set(context.Background(), key, value, 0).Err()
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig() + fmt.Sprintf(`
data "redis_string" "test" {
  key              = "%s"
  max_wait_seconds = 10
}
`, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "value", value),
				),
			},
		},
	})
}

func newRealRedisClient(t *testing.T) *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	require.NotEmpty(t, redisURL, "REDIS_URL must be set")
	opt, err := redis.ParseURL(redisURL)
	require.NoError(t, err)

	return redis.NewClient(opt)
}
