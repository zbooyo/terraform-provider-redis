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
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "Timeout for the Redis GET operation in seconds.",
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Maximum number of retries for Redis GET operation.",
			},
		},
	}
}

func dataSourceRedisStringRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*ProviderConfig)
	key := d.Get("key").(string)

	timeoutSec := d.Get("timeout").(int)
	maxRetries := d.Get("max_retries").(int)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	var val string
	var err error

	for i := 0; i < maxRetries; i++ {
		val, err = cfg.RedisClient.Get(timeoutCtx, key).Result()
		if err == nil || err == redis.Nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err == redis.Nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(key)
	_ = d.Set("key", key)
	_ = d.Set("value", val)

	return nil
}
