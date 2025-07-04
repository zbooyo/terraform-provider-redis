package redisprovider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/redis/go-redis/v9"
)

type ProviderConfig struct {
	RedisClient *redis.Client
	Timeout     time.Duration
	MaxRetries  int
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"redis_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Redis server URL, e.g. redis://localhost:6379/0",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"redis_string": resourceRedisString(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"redis_key": dataSourceRedisString(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	redisURL := d.Get("redis_url").(string)
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, diag.Errorf("Failed to parse redis_url: %s", err)
	}
	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, diag.Errorf("Failed to connect to Redis: %s", err)
	}
	return &ProviderConfig{RedisClient: client}, nil
}
