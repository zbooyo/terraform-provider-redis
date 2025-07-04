terraform {
    required_providers {
        redis = {
            source = "local/rdeavila94/redis"
            version = "~> 0.0"
        }
    }
}

provider "redis" {
    redis_url = "redis://localhost:6379/0"
}