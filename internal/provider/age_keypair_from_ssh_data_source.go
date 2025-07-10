package provider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	sshage "github.com/Mic92/ssh-to-age"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ageKeyPairFromSSHDataSource{}
	_ datasource.DataSourceWithConfigure = &ageKeyPairFromSSHDataSource{}
)

// NewageKeyPairFromSSHDataSource is a helper function to simplify the provider implementation.
func NewageKeyPairFromSSHDataSource() datasource.DataSource {
	return &ageKeyPairFromSSHDataSource{}
}

// ageKeyPairFromSSHDataSource is the data source implementation.
type ageKeyPairFromSSHDataSource struct {
}

// ageKeyPairFromSSHDataSourceModel maps the data source schema data.
type ageKeyPairFromSSHDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	SSHPrivateKey types.String `tfsdk:"ssh_private_key"`
	AgePrivateKey types.String `tfsdk:"age_private_key"`
	AgePublicKey  types.String `tfsdk:"age_public_key"`
}

// Configure adds the provider configured client to the data source.
func (d *ageKeyPairFromSSHDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	// No configuration needed for this data source
}

// Metadata returns the data source type name.
func (d *ageKeyPairFromSSHDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypair_from_ssh"
}

// Schema defines the schema for the data source.
func (d *ageKeyPairFromSSHDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Converts an Ed25519 SSH private key to an age key pair.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the data source.",
				Computed:    true,
			},
			"ssh_private_key": schema.StringAttribute{
				Description: "The Ed25519 SSH private key to convert.",
				Required:    true,
				Sensitive:   true,
			},
			"age_private_key": schema.StringAttribute{
				Description: "The converted age private key.",
				Computed:    true,
				Sensitive:   true,
			},
			"age_public_key": schema.StringAttribute{
				Description: "The converted age public key.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ageKeyPairFromSSHDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ageKeyPairFromSSHDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshPrivateKey := state.SSHPrivateKey.ValueString()

	priv, pub, err := sshage.SSHPrivateKeyToAge([]byte(sshPrivateKey), []byte(""))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting SSH Private Key to Age",
			fmt.Sprintf("Could not convert SSH private key to age: %s", err),
		)
	}

	h := sha256.New()
	h.Write([]byte(sshPrivateKey))
	id := base64.StdEncoding.EncodeToString(h.Sum(nil))

	state.AgePrivateKey = types.StringValue(*priv)
	state.AgePublicKey = types.StringValue(*pub)

	state.ID = types.StringValue(id)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
