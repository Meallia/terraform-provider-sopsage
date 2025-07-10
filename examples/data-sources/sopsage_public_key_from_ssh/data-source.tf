data "sopsage_public_key_from_ssh" "example" {
  ssh_public_key = file("~/.ssh/id_ed25519.pub")
}
