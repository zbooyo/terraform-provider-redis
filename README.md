# Terraform Provider for Redis

A Terraform provider for managing Redis keys. Useful for when Redis is used as a store for infrastructure or application configuration.
**A big thank you to [rdeavila94](https://github.com/rdeavila94) for creating the original version of this provider that this project is forked from!**

Clone of https://github.com/zbooyo/terraform-provider-redis repo with added datasource support.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0  
- [Go](https://golang.org/doc/install) >= 1.19

## Installation

### Using Terraform CLI

```bash
terraform init
```

### Manual Installation

1. Download the latest release for your platform
2. Extract the binary to your Terraform plugins directory:
  * Linux/macOS: ~/.terraform.d/plugins/registry.terraform.io/zbooyo/redis/
  * Windows: %APPDATA%\terraform.d\plugins\registry.terraform.io\zbooyo\redis\

## Usage

### Provider Configuration

```hcl
terraform {
  required_providers {
    redis = {
      source  = "zbooyo/redis"
      version = "~> 0.0"
    }
  }
}

provider "redis" {
  redis_url = "redis://localhost:6379/0"
}
```

### Resources

#### `redis_string`

Manages a Redis string key-value pair.

```hcl
resource "redis_string" "example" {
  key         = "my_key"
  value       = "my_value"
  overridable = false  # optional, defaults to false
}
```

**Arguments:**

- key (Required) - The Redis key to manage
- value (Required) - The string value to store
- overridable (Optional) - If true, allows overriding existing Redis keys. Defaults to false.

**Attributes:**

- `id` - The Redis key
- `key` - The Redis key
- `value` - The stored value
- `overridable` - Whether the key can override existing values

### Data sources 

Reads the value of a Redis string key. 
Waits until the key exists and has a non-empty value, or until a timeout is reached.

#### `redis_string`

```hcl
data "redis_string" "example" {
  key             = "my_key"
  max_wait_seconds = 300  # Optional; defaults to 300 seconds
}
```

**Arguments:**

- `key` (Required) — The Redis key to read
- `max_wait_seconds` (Optional) — Maximum time in seconds to wait for the Redis key to exist and have a non-empty value. Defaults to 300.

**Attributes:**

- `id` - The Redis key
- `key` - The Redis key
- `value` - The stored value


## Examples

### Basic Usage

```hcl
terraform {
  required_providers {
    redis = {
      source  = "zbooyo/redis"
      version = "~> 0.0"
    }
  }
}

provider "redis" {
  redis_url = "rediss://localhost:6379/0"
}

resource "redis_string" "app_config" {
  key   = "app:config:version"
  value = "1.0.0"
}

resource "redis_string" "user_session" {
  key   = "user:session:12345"
  value = "active"
}

data "redis_string" "app_version" {
  key = redis_string.app_config.key
}
```

### Overriding Existing Keys

```hcl
resource "redis_string" "existing_key" {
  key         = "existing:key"
  value       = "new_value"
  overridable = true
}
```


## Development

### Building from Source

```bash
git clone https://github.com/zbooyo/terraform-provider-redis
cd terraform-provider-redis
make build
```

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue on GitHub.
