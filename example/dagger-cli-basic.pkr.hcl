packer {
  required_plugins {
    dagger-cli = {
      source  = "github.com/hashicorp/dagger-cli"
      version = ">= 0.1.0"
    }
  }
}

source "null" "example" {
  communicator = "none"
}

build {
  sources = ["source.null.example"]

  # Basic example with a GitHub module
  provisioner "dagger-cli" {
    module = "github.com/dagger/dagger//ci?ref=main"
    call   = "test"

    args = {
      packages = "core"
    }

    log_level = "info"
  }

  # Example with local module and retries
  provisioner "dagger-cli" {
    module = "./ci"
    call   = "build"

    args = {
      platform = "linux/amd64"
      cache    = true
    }

    env = {
      BUILD_ENV = "production"
    }

    retry_count   = 2
    retry_backoff = "5s"
  }

  # Example with keyed caching
  provisioner "dagger-cli" {
    module = "./infrastructure"
    call   = "provision"

    args = {
      region = "us-east-1"
    }

    cache_mode = "keyed"
    cache_key  = "infrastructure-base"
  }
}