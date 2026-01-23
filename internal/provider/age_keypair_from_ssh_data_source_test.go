package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAgeKeyPairFromSSHDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "sopsage_keypair_from_ssh" "test" {
					  ssh_private_key = <<EOT%sEOT
					}`, sshPrivkey),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.sopsage_keypair_from_ssh.test",
						tfjsonpath.New("age_public_key"),
						knownvalue.StringExact(agePubkey),
					),
					statecheck.ExpectKnownValue(
						"data.sopsage_keypair_from_ssh.test",
						tfjsonpath.New("age_private_key"),
						knownvalue.StringExact(agePrivkey),
					),
				},
			},
		},
	})
}
