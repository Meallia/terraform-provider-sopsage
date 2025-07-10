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
	_ datasource.DataSource              = &agePublicKeyFromSSHDataSource{}
	_ datasource.DataSourceWithConfigure = &agePublicKeyFromSSHDataSource{}
)

// NewAgePublicKeyFromSSHDataSource is a helper function to simplify the provider implementation.
func NewAgePublicKeyFromSSHDataSource() datasource.DataSource {
	return &agePublicKeyFromSSHDataSource{}
}

// agePublicKeyFromSSHDataSource is the data source implementation.
type agePublicKeyFromSSHDataSource struct {
}

// agePublicKeyFromSSHDataSourceModel maps the data source schema data.
type agePublicKeyFromSSHDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	SSHPublicKey types.String `tfsdk:"ssh_public_key"`
	AgePublicKey types.String `tfsdk:"age_public_key"`
}

// Configure adds the provider configured client to the data source.
func (d *agePublicKeyFromSSHDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	// No configuration needed for this data source
}

// Metadata returns the data source type name.
func (d *agePublicKeyFromSSHDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_key_from_ssh"
}

// Schema defines the schema for the data source.
func (d *agePublicKeyFromSSHDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Extracts an age public key from an Ed25519 SSH public key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the data source.",
				Computed:    true,
			},
			"ssh_public_key": schema.StringAttribute{
				Description: "The Ed25519 SSH public key to extract the age public key from.",
				Required:    true,
				Sensitive:   false,
			},
			"age_public_key": schema.StringAttribute{
				Description: "The extracted age public key.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *agePublicKeyFromSSHDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state agePublicKeyFromSSHDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshPublicKey := state.SSHPublicKey.ValueString()

	agePublicKeyPtr, err := sshage.SSHPublicKeyToAge([]byte(sshPublicKey))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting SSH Public Key to Age",
			fmt.Sprintf("Could not convert SSH public key to age: %s", err),
		)
		return
	}

	h := sha256.New()
	h.Write([]byte(sshPublicKey))
	id := base64.StdEncoding.EncodeToString(h.Sum(nil))

	state.AgePublicKey = types.StringValue(*agePublicKeyPtr)

	state.ID = types.StringValue(id)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
