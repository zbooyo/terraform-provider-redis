package redisprovider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/redis/go-redis/v9"
)

func dataSourceRedisString() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRedisStringRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Redis key to read.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value stored at the Redis key.",
			},
			"max_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Maximum time in seconds to wait for the Redis key to exist and have a non-empty value.",
			},
		},
	}
}

func dataSourceRedisStringRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*ProviderConfig)
	key := d.Get("key").(string)

	maxWait := d.Get("max_wait_seconds").(int)
	deadline := time.Now().Add(time.Duration(maxWait) * time.Second)

	var val string
	var err error

	for time.Now().Before(deadline) {
		val, err = cfg.RedisClient.Get(ctx, key).Result()

		if err == redis.Nil || val == "" {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(key)
		_ = d.Set("key", key)
		_ = d.Set("value", val)
		return nil
	}

	d.SetId("")
	return nil
}
