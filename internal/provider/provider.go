// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thecodedcloud/terraform-provider-discord/internal/provider/channel"
	"github.com/thecodedcloud/terraform-provider-discord/internal/provider/guild"
	"github.com/thecodedcloud/terraform-provider-discord/internal/provider/member"
	"github.com/thecodedcloud/terraform-provider-discord/internal/provider/permissions"
	"github.com/thecodedcloud/terraform-provider-discord/internal/provider/role"
	"github.com/thecodedcloud/terraform-provider-discord/internal/provider/role_members"
	"github.com/thecodedcloud/terraform-provider-discord/internal/provider/webhook"
)

// Ensure DiscordProvider satisfies various provider interfaces.
var _ provider.Provider = &DiscordProvider{}
var _ provider.ProviderWithFunctions = &DiscordProvider{}

// DiscordProvider defines the provider implementation.
type DiscordProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DiscordProviderModel describes the provider data model.
type DiscordProviderModel struct {
	AccessToken    types.String `tfsdk:"access_token"`
	OAuth2ClientId types.String `tfsdk:"oauth2_client_id"`
}

// New is a helper function to simplfy provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DiscordProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name and version.
func (p *DiscordProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "discord"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *DiscordProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"access_token": schema.StringAttribute{
				MarkdownDescription: "The bot access token for Discord.",
				Optional:            true,
				Sensitive:           true,
			},
			"oauth2_client_id": schema.StringAttribute{
				MarkdownDescription: "The OAuth2 client ID for Discord.",
				Optional:            true,
			},
		},
	}
}

// Configure prepares a Discord API client for use by data sources and resources.
func (p *DiscordProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from the configuration.
	var config DiscordProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// All configuration values must be known values.
	if config.AccessToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_token"),
			"Unknown Discord access token",
			"The provider cannot create the Discord client as there is an unknown configuration value for the Discord access token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DISCORD_ACCESS_TOKEN environment variable.",
		)
	}

	if config.OAuth2ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("oauth2_client_id"),
			"Unknown Discord OAuth2 client ID",
			"The provider cannot create the Discord client as there is an unknown configuration value for the Discord OAuth2 client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DISCORD_OAUTH2_CLIENT_ID environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override with Terraform configuration value if set.
	access_token := os.Getenv("DISCORD_ACCESS_TOKEN")
	oauth2_client_id := os.Getenv("DISCORD_OAUTH2_CLIENT_ID")

	if !config.AccessToken.IsNull() {
		access_token = config.AccessToken.ValueString()
	}

	if !config.OAuth2ClientId.IsNull() {
		oauth2_client_id = config.OAuth2ClientId.ValueString()
	}

	// If any of the required configurations are missing, return errors with specific guidance.
	if access_token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_token"),
			"Missing Discord access token",
			"The provider cannot create the Discord client as there is a missing or empty value for the Discord access token. "+
				"Set the access_token value in the configuration or use the DISCORD_ACCESS_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if oauth2_client_id == "" {
		// todo, this is an optional value, so we should not return an error here.
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Discord client using the configuration values.
	client, err := discordgo.New("Bot " + access_token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Discord client",
			"An unexpected error occurred when creating the Discord client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Discord client Error: "+err.Error(),
		)
		return
	}

	// Make the client available to data sources and resources type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources defines the list of resources implemented by the provider.
func (p *DiscordProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		channel.NewRoleResource,
		role.NewRoleResource,
		permissions.NewPermissionsResource,
		webhook.NewWebhookResource,
		role_members.NewRoleMembersResource,
	}
}

// DataSources defines the list of data sources implemented by the provider.
func (p *DiscordProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		guild.NewGuildDataSource,
		channel.NewChannelDataSource,
		role.NewRoleDataSource,
		permissions.NewPermissionsDataSource,
		webhook.NewWebhookDataSource,
		member.NewMemberDataSource,
		role_members.NewRoleMembersDataSource,
	}
}

// Functions defines the list of functions implemented by the provider.
func (p *DiscordProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}
