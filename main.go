package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"terraform-provider-redis/redisprovider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: redisprovider.Provider,
	})
} 