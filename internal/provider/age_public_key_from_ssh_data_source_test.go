package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAgePublicKeyFromSSHDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "sopsage_public_key_from_ssh" "test" {
					  ssh_public_key = "%s"
					}`, sshPubkey),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.sopsage_public_key_from_ssh.test",
						tfjsonpath.New("age_public_key"),
						knownvalue.StringExact(agePubkey),
					),
				},
			},
		},
	})
}
