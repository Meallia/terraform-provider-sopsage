resource "sopsage_encrypted_data" "test" {
  format  = "yaml"
  content = yamlencode({ foo = "bar" })
  age_public_keys = [
    "age1c2cnfzjfeswsydufz4tcrs46zqpmnz9t3dwz5uaef5yl3qnzaptqy88dl7"
  ]
}
