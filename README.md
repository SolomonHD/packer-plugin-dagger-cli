# Packer Plugin for Dagger CLI

This plugin enables integration between [Packer](https://www.packer.io) and [Dagger](https://dagger.io), allowing you to execute Dagger CLI commands during image provisioning. This brings container-based reproducible workflows to Packer builds.

## Features

- **dagger-cli provisioner**: Execute Dagger CLI commands during provisioning
- **Retry mechanisms**: Configurable retry logic with exponential backoff
- **Caching control**: Support for Dagger's caching modes (auto, disabled, keyed)
- **Structured logging**: Detailed logging with configurable verbosity levels
- **Full argument support**: Pass arguments, environment variables, and configuration
- **Local and remote modules**: Support for both local and remote Dagger modules

## Installation

### Using `packer init` (Recommended)

Add the plugin to your Packer template:

```hcl
packer {
  required_plugins {
    dagger-cli = {
      source  = "github.com/hashicorp/dagger-cli"
      version = ">= 0.1.0"
    }
  }
}
```

Then run:

```shell
packer init .
```

### Manual Installation

Download the latest release from the [releases page](https://github.com/hashicorp/packer-plugin-dagger-cli/releases) and install using:

```shell
packer plugins install --path packer-plugin-dagger-cli github.com/hashicorp/dagger-cli
```

## Usage

### Basic Example

```hcl
provisioner "dagger-cli" {
  module = "github.com/acme/infra//ci?ref=v1.2.3"
  call   = "configure_system"

  args = {
    region = "us-east-1"
    fast   = true
  }
}
```

### With Retries

```hcl
provisioner "dagger-cli" {
  module = "./ci"
  call   = "deploy"

  retry_count   = 3
  retry_backoff = "10s"

  args = {
    environment = "production"
  }
}
```

### With Keyed Caching

```hcl
provisioner "dagger-cli" {
  module = "./infrastructure"
  call   = "provision"

  cache_mode = "keyed"
  cache_key  = "ubuntu-22.04-base"

  args = {
    os_version = "22.04"
  }
}
```

See the [documentation](docs/provisioners/dagger-cli.mdx) for complete configuration options and examples.

## Requirements

- **Packer**: >= 1.10.2
- **Go**: >= 1.23 (for building from source)
- **Dagger CLI**: Must be installed and available in PATH
- **packer-plugin-sdk**: >= 0.6.1

## Build from Source

1. Clone this repository:

   ```shell
   git clone https://github.com/hashicorp/packer-plugin-dagger-cli.git
   cd packer-plugin-dagger-cli
   ```

2. Build the plugin:

   ```shell
   VERSION=$(cat version/VERSION)
   go build -ldflags="-X github.com/hashicorp/packer-plugin-dagger-cli/version.Version=$VERSION -X github.com/hashicorp/packer-plugin-dagger-cli/version.VersionPrerelease=" -o packer-plugin-dagger-cli
   ```

3. Install the plugin:

   ```shell
   packer plugins install --path packer-plugin-dagger-cli github.com/hashicorp/dagger-cli
   ```

### Quick Build (Unix-like systems)

```shell
make dev
```

## Development

### Running Tests

Run unit tests:

```shell
go test ./provisioner/dagger-cli/ -v
```

Run all tests:

```shell
go test -race -count 1 ./... -timeout=3m
```

### Running Acceptance Tests

```shell
make testacc
```

Or manually:

```shell
PACKER_ACC=1 go test -count 1 -v ./... -timeout=120m
```

### Code Generation

Generate HCL2 specs:

```shell
go generate ./...
```

Or:

```shell
make generate
```

## Documentation

- [Provisioner Documentation](docs/provisioners/dagger-cli.mdx) - Complete configuration reference
- [Examples](example/) - Working example templates
- [Packer Plugins](https://www.packer.io/docs/plugins) - General Packer plugin documentation
- [Dagger Documentation](https://docs.dagger.io) - Dagger CLI and module documentation

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## License

MPL-2.0 License - see [LICENSE](LICENSE) for details.

## Support

- [GitHub Issues](https://github.com/hashicorp/packer-plugin-dagger-cli/issues) - Bug reports and feature requests
- [Packer Community Forum](https://discuss.hashicorp.com/c/packer) - General questions and discussions
- [Dagger Discord](https://discord.gg/dagger-io) - Dagger-specific questions
