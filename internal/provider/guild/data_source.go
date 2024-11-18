package guild

import (
	"context"
	"fmt"

	"github.com/JustARecord/go-discordutils/base/guild"
	discord "github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
)

// NewGuildDataSource is a helper function to simplify the provider implementation.
func NewGuildDataSource() datasource.DataSource {
	return &GuildDataSource{}
}

// Metadata returns the data source type name.
func (d *GuildDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceMetadataName
}

// Schema defines the schema for the data source.
func (d *GuildDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the guild.",
				Optional:    true,
				Computed:    true,
			},
			"icon": schema.StringAttribute{
				Description: "The hash of the guild's icon.",
				Computed:    true,
			},
			"splash": schema.StringAttribute{
				Description: "The hash of the guild's splash.",
				Computed:    true,
			},
			"discovery_splash": schema.StringAttribute{
				Description: "The hash of the guild's discovery splash.",
				Computed:    true,
			},
			"owner": schema.BoolAttribute{
				Description: "If we are the owner of the guild.",
				Computed:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "The user ID of the owner of the guild.",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "The voice region of the guild.",
				Computed:    true,
			},
			"afk_channel_id": schema.StringAttribute{
				Description: "The ID of the AFK voice channel.",
				Computed:    true,
			},
			"afk_timeout": schema.Int32Attribute{
				Description: "The timeout, in seconds, before a user is considered AFK in voice.",
				Computed:    true,
			},
			"widget_enabled": schema.BoolAttribute{
				Description: "Whether or not the Server Widget is enabled.",
				Computed:    true,
			},
			"widget_channel_id": schema.StringAttribute{
				Description: "The Channel ID for the Server Widget.",
				Computed:    true,
			},
			"verification_level": schema.StringAttribute{
				Description: "The verification level required for the guild.",
				Computed:    true,
			},
			"default_message_notifications": schema.StringAttribute{
				Description: "The default message notification setting for the guild.",
				Computed:    true,
			},
			"explicit_content_filter": schema.StringAttribute{
				Description: "The explicit content filter level.",
				Computed:    true,
			},
			"features": schema.ListAttribute{
				Description: "Enabled guild features.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"mfa_level": schema.StringAttribute{
				Description: "Required MFA level for the guild.",
				Computed:    true,
			},
			"application_id": schema.StringAttribute{
				Description: "The application id of the guild if bot created.",
				Computed:    true,
			},
			"system_channel_id": schema.StringAttribute{
				Description: "The Channel ID to which system messages are sent (eg join and leave messages).",
				Computed:    true,
			},
			"system_channel_flags": schema.ListAttribute{
				Description: "The System channel flags.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"rules_channel_id": schema.StringAttribute{
				Description: "The ID of the rules channel ID, used for rules.",
				Computed:    true,
			},
			"max_presences": schema.Int32Attribute{
				Description: "The maximum number of presences for the guild.",
				Computed:    true,
			},
			"max_members": schema.Int32Attribute{
				Description: "The maximum number of members for the guild.",
				Computed:    true,
			},
			"vanity_url_code": schema.StringAttribute{
				Description: "The vanity url code for the guild.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the guild.",
				Computed:    true,
			},
			"banner": schema.StringAttribute{
				Description: "The hash of the guild's banner.",
				Computed:    true,
			},
			"premium_tier": schema.StringAttribute{
				Description: "The premium tier of the guild (Server Boost Level).",
				Computed:    true,
			},
			"premium_subscription_count": schema.Int32Attribute{
				Description: "The total number of users currently boosting this server.",
				Computed:    true,
			},
			"preferred_locale": schema.StringAttribute{
				Description: "The preferred locale of a guild with the 'PUBLIC' feature; used in server discovery and notices from Discord; defaults to 'en-US'.",
				Computed:    true,
			},
			"public_updates_channel_id": schema.StringAttribute{
				Description: "The id of the channel where admins and moderators of guilds with the 'PUBLIC' feature receive notices from Discord.",
				Computed:    true,
			},
			"max_video_channel_users": schema.Int32Attribute{
				Description: "The maximum amount of users in a video channel.",
				Computed:    true,
			},
			"max_stage_video_channel_users": schema.Int32Attribute{
				Description: "The maximum amount of users in a stage video channel.",
				Computed:    true,
			},
			"approximate_member_count": schema.Int32Attribute{
				Description: "Approximate number of members in this guild, returned from the GET /guild/<id> endpoint when with_counts is true.",
				Computed:    true,
			},
			"approximate_presence_count": schema.Int32Attribute{
				Description: "Approximate number of non-offline members in this guild, returned from the GET /guild/<id> endpoint when with_counts is true.",
				Computed:    true,
			},
			"nsfw_level": schema.StringAttribute{
				Description: "The NSFW Level of the guild.",
				Computed:    true,
			},
			"premium_progress_bar_enabled": schema.BoolAttribute{
				Description: "Whether the guild has the boost progress bar enabled.",
				Computed:    true,
			},
			"safety_alerts_channel_id": schema.StringAttribute{
				Description: "The id of the channel where admins and moderators of Community guilds receive safety alerts from Discord.",
				Computed:    true,
			},
			"permissions": schema.Int64Attribute{
				Description: "Permissions of our user.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *GuildDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var state, provided GuildDataSourceModel

	// Read the configuration data into the provided struct.
	diags := req.Config.Get(ctx, &provided)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"name?": provided.Name,
		"id?":   provided.ID,
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

	var result *discordgo.Guild
	var err error

	// Fetch data from the Discord client
	if id != "" {
		result, err = guild.FetchByID(ctx, d.client, id)
	} else if name != "" {
		result, err = guild.FetchByName(ctx, d.client, name)
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

	systemChannelFlags := discord.ListStringify(result.SystemChannelFlags)
	features := discord.ListStringify(result.Features)

	systemChannelFlagsList, diags := common.ToListType[string, basetypes.StringType](systemChannelFlags)
	resp.Diagnostics.Append(diags...)

	featuresList, diags := common.ToListType[string, basetypes.StringType](features)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Map the result data to the state.
	state = GuildDataSourceModel{
		ID:                          types.StringValue(result.ID),
		Name:                        types.StringValue(result.Name),
		Icon:                        types.StringValue(result.Icon),
		Splash:                      types.StringValue(result.Splash),
		DiscoverySplash:             types.StringValue(result.DiscoverySplash),
		Owner:                       types.BoolValue(result.Owner),
		OwnerID:                     types.StringValue(result.OwnerID),
		Region:                      types.StringValue(result.Region),
		AfkChannelID:                types.StringValue(result.AfkChannelID),
		AfkTimeout:                  types.Int32Value(int32(result.AfkTimeout)),
		WidgetEnabled:               types.BoolValue(result.WidgetEnabled),
		WidgetChannelID:             types.StringValue(result.WidgetChannelID),
		VerificationLevel:           types.StringValue(discord.Stringify(result.VerificationLevel)),
		DefaultMessageNotifications: types.StringValue(discord.Stringify(result.DefaultMessageNotifications)),
		ExplicitContentFilter:       types.StringValue(discord.Stringify(result.ExplicitContentFilter)),
		Features:                    featuresList,
		MfaLevel:                    types.StringValue(discord.Stringify(result.MfaLevel)),
		ApplicationID:               types.StringValue(result.ApplicationID),
		SystemChannelID:             types.StringValue(result.SystemChannelID),
		SystemChannelFlags:          systemChannelFlagsList,
		RulesChannelID:              types.StringValue(result.RulesChannelID),
		MaxPresences:                types.Int32Value(int32(result.MaxPresences)),
		MaxMembers:                  types.Int32Value(int32(result.MaxMembers)),
		VanityURLCode:               types.StringValue(result.VanityURLCode),
		Description:                 types.StringValue(result.Description),
		Banner:                      types.StringValue(result.Banner),
		PremiumTier:                 types.StringValue(discord.Stringify(result.PremiumTier)),
		PremiumSubscriptionCount:    types.Int32Value(int32(result.PremiumSubscriptionCount)),
		PreferredLocale:             types.StringValue(result.PreferredLocale),
		PublicUpdatesChannelID:      types.StringValue(result.PublicUpdatesChannelID),
		MaxVideoChannelUsers:        types.Int32Value(int32(result.MaxVideoChannelUsers)),
		// MaxStageVideoChannelUsers:   types.Int32Value(int32(guild.MaxStageVideoChannelUsers)),
		ApproximateMemberCount:   types.Int32Value(int32(result.ApproximateMemberCount)),
		ApproximatePresenceCount: types.Int32Value(int32(result.ApproximatePresenceCount)),
		NSFWLevel:                types.StringValue(discord.Stringify(result.NSFWLevel)),
		// PremiumProgressBarEnabled: types.BoolValue(guild.PremiumProgressBarEnabled),
		// SafetyAlertsChannelID:     types.StringValue(guild.SafetyAlertsChannelID),
		Permissions: types.Int64Value(result.Permissions),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *GuildDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
