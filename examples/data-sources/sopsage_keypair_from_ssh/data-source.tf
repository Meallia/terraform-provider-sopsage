data "sopsage_keypair_from_ssh" "example" {
  ssh_private_key = file("~/.ssh/id_ed25519")
}
