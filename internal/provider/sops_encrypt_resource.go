package provider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &sopsEncryptResource{}
	_ resource.ResourceWithConfigure = &sopsEncryptResource{}
)

// NewSopsEncryptResource is a helper function to simplify the provider implementation.
func NewSopsEncryptResource() resource.Resource {
	return &sopsEncryptResource{}
}

// sopsEncryptResource is the resource implementation.
type sopsEncryptResource struct {
}

// sopsEncryptResourceModel maps the resource schema data.
type sopsEncryptResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Content       types.String `tfsdk:"content"`
	Format        types.String `tfsdk:"format"`
	AgePublicKeys types.List   `tfsdk:"age_public_keys"`
	Encrypted     types.String `tfsdk:"encrypted"`
}

// Configure adds the provider configured client to the resource.
func (r *sopsEncryptResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	// No configuration needed for this resource
}

// Metadata returns the resource type name.
func (r *sopsEncryptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_encrypted_data"
}

// Schema defines the schema for the resource.
func (r *sopsEncryptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Encrypts content using SOPS with age encryption.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"content": schema.StringAttribute{
				Description: "The content to encrypt.",
				Required:    true,
			},
			"format": schema.StringAttribute{
				Description: "The format of the content (json, yaml, etc.).",
				Required:    true,
			},
			"age_public_keys": schema.ListAttribute{
				Description: "List of age public keys to encrypt with.",
				Required:    true,
				ElementType: types.StringType,
			},
			"encrypted": schema.StringAttribute{
				Description: "The encrypted content in SOPS format.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

// Create creates a new encrypted content.
func (r *sopsEncryptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sopsEncryptResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the content and format
	content := plan.Content.ValueString()
	format := plan.Format.ValueString()

	// Get the age public keys
	var agePublicKeys []string
	diags = plan.AgePublicKeys.ElementsAs(ctx, &agePublicKeys, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Encrypt the content
	encrypted, err := SopsEncryptDataFromAgeKeys(content, format, agePublicKeys)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Encrypting Content",
			fmt.Sprintf("Could not encrypt content: %s", err),
		)
		return
	}

	// Generate a unique ID based on the content and keys
	h := sha256.New()
	h.Write([]byte(content))
	for _, key := range agePublicKeys {
		h.Write([]byte(key))
	}
	id := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Set resource ID
	plan.ID = types.StringValue(id)
	// Set encrypted content
	plan.Encrypted = types.StringValue(encrypted)

	// Set state to computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *sopsEncryptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state sopsEncryptResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sopsEncryptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan sopsEncryptResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state sopsEncryptResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if content or keys have changed
	contentChanged := plan.Content.ValueString() != state.Content.ValueString()
	formatChanged := plan.Format.ValueString() != state.Format.ValueString()

	var planKeys []string
	diags = plan.AgePublicKeys.ElementsAs(ctx, &planKeys, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateKeys []string
	diags = state.AgePublicKeys.ElementsAs(ctx, &stateKeys, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	keysChanged := len(planKeys) != len(stateKeys)
	if !keysChanged {
		for i, key := range planKeys {
			if key != stateKeys[i] {
				keysChanged = true
				break
			}
		}
	}

	// Only re-encrypt if content, format, or keys have changed
	if contentChanged || formatChanged || keysChanged {
		// Get the content and format
		content := plan.Content.ValueString()
		format := plan.Format.ValueString()

		// Encrypt the content
		encrypted, err := SopsEncryptDataFromAgeKeys(content, format, planKeys)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Encrypting Content",
				fmt.Sprintf("Could not encrypt content: %s", err),
			)
			return
		}

		// Generate a new unique ID
		h := sha256.New()
		h.Write([]byte(content))
		for _, key := range planKeys {
			h.Write([]byte(key))
		}
		id := base64.StdEncoding.EncodeToString(h.Sum(nil))

		// Update ID and encrypted content
		plan.ID = types.StringValue(id)
		plan.Encrypted = types.StringValue(encrypted)
	} else {
		// Keep the existing encrypted content and ID
		plan.ID = state.ID
		plan.Encrypted = state.Encrypted
	}

	// Set state to computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sopsEncryptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Encrypted content doesn't have any external resources to clean up
	// The state will be removed by Terraform automatically
}
