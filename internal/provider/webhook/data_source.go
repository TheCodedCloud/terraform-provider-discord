package webhook

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
	"github.com/justarecord/terraform-provider-discord/internal/provider/discord"
)

// NewWebhookDataSource is a helper function to simplify the provider implementation.
func NewWebhookDataSource() datasource.DataSource {
	return &WebhookDataSource{}
}

// Metadata returns the data source type name.
func (d *WebhookDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceMetadataName
}

// Schema defines the schema for the data source.
func (d *WebhookDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the webhook.",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the webhook, either 'INCOMING', 'CHANNEL_FOLLOWER', or 'APPLICATION'.",
				Computed:    true,
			},
			"guild_id": schema.StringAttribute{
				Description: "The guild ID this webhook is for, if any.",
				Optional:    true,
				Computed:    true,
			},
			"channel_id": schema.StringAttribute{
				Description: "The channel ID this webhook is for, if any.",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The default name of the webhook.",
				Optional:    true,
				Computed:    true,
			},
			"avatar": schema.StringAttribute{
				Description: "The default user avatar hash of the webhook.",
				Computed:    true,
				Sensitive:   true, // Marked as sensitive to prevent log spam with large base64 encoded images.
			},
			"token": schema.StringAttribute{
				Description: "The secure token of the webhook (returned for Incoming Webhooks).",
				Computed:    true,
				Sensitive:   true,
			},
			"application_id": schema.StringAttribute{
				Description: "The Bot/OAuth2 application that created this webhook.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *WebhookDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var state, provided WebhookDataSourceModel

	// Read the configuration data into the provided struct.
	diags := req.Config.Get(ctx, &provided)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"id?":         provided.ID,
		"name?":       provided.Name,
		"guild_id?":   provided.GuildID,
		"channel_id?": provided.ChannelID,
	}

	// Check for null/unknown values
	common.CheckNonNull(required, resp.Diagnostics, datasourceMetadataName, datasourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required values
	common.CheckRequired(ctx, required, resp.Diagnostics, datasourceMetadataName, datasourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	id := provided.ID.ValueString()
	name := provided.Name.ValueString()
	guild_id := provided.GuildID.ValueString()
	channel_id := provided.ChannelID.ValueString()

	var result *discordgo.Webhook
	var err error

	// Fetch data from the Discord client
	if id != "" {
		result, err = discord.FetchWebhookByID(ctx, d.client, id)
	} else if name != "" {
		if channel_id != "" {
			result, err = discord.FetchChannelWebhookByName(ctx, d.client, channel_id, name)
		} else if guild_id != "" {
			result, err = discord.FetchGuildWebhookByName(ctx, d.client, guild_id, name)
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Invalid %s Configuration", datasourceMetadataType),
				fmt.Sprintf("Either the guild_id or the channel_id must be set if name is set for the %s %s.", datasourceMetadataName, datasourceMetadataType),
			)
		}
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s Configuration", datasourceMetadataType),
			fmt.Sprintf("Either the id or the name must be set for the %s %s.", datasourceMetadataName, datasourceMetadataType),
		)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", datasourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Map the result data to the state.
	state = WebhookDataSourceModel{
		ID:            types.StringValue(result.ID),
		Type:          types.StringValue(discord.Stringify(result.Type)),
		GuildID:       types.StringValue(result.GuildID),
		ChannelID:     types.StringValue(result.ChannelID),
		Name:          types.StringValue(result.Name),
		Avatar:        types.StringValue(result.Avatar),
		Token:         types.StringValue(result.Token),
		ApplicationID: types.StringValue(result.ApplicationID),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *WebhookDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	// Get the client from the provider data.
	client, ok := req.ProviderData.(*discordgo.Session)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unexpected %s Configure Type", datasourceMetadataType),
			fmt.Sprintf("Expected *discordgo.Session, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
