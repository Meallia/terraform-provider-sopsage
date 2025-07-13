# Terraform SOPS Age Provider

This provider allows you to work with [SOPS](https://github.com/mozilla/sops)
and [Age](https://age-encryption.org/) encryption directly in your Terraform configurations.

The provider is available on both
the [Terraform Registry](https://registry.terraform.io/providers/Meallia/sopsage/latest/docs)
and the [OpenTofu Registry](https://search.opentofu.org/provider/meallia/sopsage/latest).

## Features

- Generate Age key pairs
- Convert Ed25519 SSH keys to Age keys
- Encrypt content using SOPS with Age encryption
