package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestSopsEncryptResourceResource(t *testing.T) {
	config := fmt.Sprintf(`
					data "sopsage_public_key_from_ssh" "test" {
					  ssh_public_key = "%s"
					}

					resource "sopsage_encrypted_data" "test" {
                      format = "yaml"
                      content = yamlencode({foo = "bar"})
					  encrypted_regex = "^f.*"
					  age_public_keys = [
						data.sopsage_public_key_from_ssh.test.age_public_key,
						data.sopsage_public_key_from_ssh.test.age_public_key
					  ]
					}
				`, sshPubkey)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
