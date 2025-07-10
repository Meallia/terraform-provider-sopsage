# Terraform SOPS Age Provider

This Terraform provider allows you to work with [SOPS](https://github.com/mozilla/sops) and [Age](https://age-encryption.org/) encryption directly in your Terraform configurations.

## Features

- Generate Age key pairs
- Convert Ed25519 SSH keys to Age keys
- Encrypt content using SOPS with Age encryption

## Installation

```hcl
terraform {
  required_providers {
    sopsage = {
      source = "meallia/sopsage"
      version = "0.1.0"
    }
  }
}

provider "sopsage" {}
```

## Resources and Data Sources

### `sopsage_keypair` Resource

Generates a new Age key pair.

```hcl
resource "sopsage_keypair" "example" {
}

output "public_key" {
  value = sopsage_keypair.example.public_key
}

output "private_key" {
  value     = sopsage_keypair.example.private_key
  sensitive = true
}
```

### `sopsage_keypair_from_ssh` Data Source

Converts an Ed25519 SSH private key to an Age key pair.

```hcl
data "sopsage_keypair_from_ssh" "example" {
  ssh_private_key = file("~/.ssh/id_ed25519")
}

output "age_public_key" {
  value = data.sopsage_keypair_from_ssh.example.age_public_key
}

output "age_private_key" {
  value     = data.sopsage_keypair_from_ssh.example.age_private_key
  sensitive = true
}
```

### `sopsage_public_key_from_ssh` Data Source

Extracts an Age public key from an Ed25519 SSH public key.

```hcl
data "sopsage_public_key_from_ssh" "example" {
  ssh_public_key = file("~/.ssh/id_ed25519.pub")
}

output "age_public_key" {
  value = data.sopsage_public_key_from_ssh.example.age_public_key
}
```

### `sopsage_encrypted_data` Resource

Encrypts content using SOPS with Age encryption.

```hcl
resource "sopsage_keypair" "example" {
}

resource "sopsage_encrypted_data" "example" {
  content = jsonencode({
    username = "admin"
    password = "supersecret"
  })
  format = "json"
  age_public_keys = [
    sopsage_keypair.example.public_key
  ]
}

output "encrypted_content" {
  value     = sopsage_encrypted_data.example.encrypted
  sensitive = true
}
```

## Requirements

- SOPS CLI must be installed and available in the PATH
- Terraform 0.13+

## Development

### Building the Provider

1. Clone the repository
2. Build the provider using `go build`

```bash
go build -o terraform-provider-sopsage
```

### Testing the Provider

```bash
go test ./...
```

## License

MIT
