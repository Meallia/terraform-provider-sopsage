package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &SopsAgeProvider{}
)

// SopsAgeProvider is the provider implementation.
type SopsAgeProvider struct {
	// version is set to the provider version on release.
	version string
}

// New creates a new provider instance.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SopsAgeProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *SopsAgeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sopsage"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *SopsAgeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with SOPS and Age encryption.",
		Attributes:  map[string]schema.Attribute{},
	}
}

// Configure prepares a SOPS Age API client for data sources and resources.
func (p *SopsAgeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// No configuration needed for this provider
}

// DataSources defines the data sources implemented in the provider.
func (p *SopsAgeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewageKeyPairFromSSHDataSource,
		NewAgePublicKeyFromSSHDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *SopsAgeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSopsEncryptResource,
	}
}
