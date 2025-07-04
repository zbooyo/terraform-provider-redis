package redisprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/redis/go-redis/v9"
)

func resourceRedisString() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedisStringCreate,
		ReadContext:   resourceRedisStringRead,
		UpdateContext: resourceRedisStringUpdate,
		DeleteContext: resourceRedisStringDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRedisStringImport,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"overridable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, allows overriding existing Redis keys. If false, creation will fail if the key already exists.",
			},
		},
	}
}

func resourceRedisStringCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*ProviderConfig)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	overridable := d.Get("overridable").(bool)

	exists, err := cfg.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return diag.FromErr(err)
	}
	if exists > 0 && !overridable {
		return diag.Errorf("redis key '%s' already exists. Set overridable = true to allow overriding existing keys", key)
	}
	if err := cfg.RedisClient.Set(ctx, key, value, 0).Err(); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(key)
	return resourceRedisStringRead(ctx, d, m)
}

func resourceRedisStringRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*ProviderConfig)
	key := d.Id()
	val, err := cfg.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		d.SetId("")
		return nil
	} else if err != nil {
		return diag.FromErr(err)
	}
	d.Set("key", key)
	d.Set("value", val)
	return nil
}

func resourceRedisStringUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*ProviderConfig)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	if err := cfg.RedisClient.Set(ctx, key, value, 0).Err(); err != nil {
		return diag.FromErr(err)
	}
	return resourceRedisStringRead(ctx, d, m)
}

func resourceRedisStringDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*ProviderConfig)
	key := d.Id()
	if err := cfg.RedisClient.Del(ctx, key).Err(); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceRedisStringImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	key := d.Id()
	cfg := m.(*ProviderConfig)
	val, err := cfg.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("redis key '%s' not found", key)
	} else if err != nil {
		return nil, err
	}
	d.Set("key", key)
	d.Set("value", val)
	return []*schema.ResourceData{d}, nil
}
