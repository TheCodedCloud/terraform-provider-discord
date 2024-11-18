package channel

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
	"github.com/justarecord/terraform-provider-discord/internal/provider/discord"
)

// NewChannelDataSource is a helper function to simplify the provider implementation.
func NewChannelDataSource() datasource.DataSource {
	return &ChannelDataSource{}
}

// Metadata returns the data source type name.
func (d *ChannelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceMetadataName
}

// Schema defines the schema for the data source.
func (d *ChannelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the channel.",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the channel.",
				Computed:    true,
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Required:    true,
			},
			"position": schema.Int32Attribute{
				Description: "The position of the channel.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the channel.",
				Optional:    true,
				Computed:    true,
			},
			"topic": schema.StringAttribute{
				Description: "The topic of the channel.",
				Computed:    true,
			},
			"nsfw": schema.BoolAttribute{
				Description: "Whether the channel is marked as NSFW.",
				Computed:    true,
			},
			"bitrate": schema.Int32Attribute{
				Description: "The bitrate of the channel, if it is a voice channel.",
				Computed:    true,
			},
			"user_limit": schema.Int32Attribute{
				Description: "The user limit of the voice channel.",
				Computed:    true,
			},
			"rate_limit_per_user": schema.Int32Attribute{
				Description: "Amount of seconds a user has to wait before sending another message or creating another thread (0-21600)",
				Computed:    true,
			},
			"icon": schema.StringAttribute{
				Description: "Icon of the group DM channel.",
				Computed:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "ID of the creator of the group DM or thread",
				Computed:    true,
			},
			"application_id": schema.StringAttribute{
				Description: "ApplicationID of the DM creator Zeroed if guild channel or not a bot user",
				Computed:    true,
			},
			// "managed": schema.BoolAttribute{
			// 	Description: "Whether the role is managed by an integration.",
			// 	Computed:    true,
			// },
			"parent_id": schema.StringAttribute{
				Description: "The ID of the parent category for a channel.",
				Computed:    true,
			},
			"children": schema.ListAttribute{
				Description: "The IDs of the child channels of the category, if the channel is a category.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"last_pin_timestamp": schema.StringAttribute{
				Description: "The timestamp of the last pinned message in the channel.",
				Computed:    true,
			},
			// "thread_metadata": schema.SingleNestedAttribute{
			// 	Description: "Thread metadata for the channel.",
			// 	Computed:    true,
			// 	Attributes:  ThreadMetadataSchema,
			// },
			// "thread_member": schema.ListAttribute{
			// 	Description: "Thread member metadata for the channel.",
			// 	Computed:    true,
			// 	CustomType:  ThreadMember,
			// },
			"flags": schema.ListAttribute{
				Description: "Channel flags.",
				Computed:    true,
				ElementType: types.StringType,
			},
			// "available_tags": schema.ListNestedAttribute{
			// 	Description:  "The set of tags that can be used in a forum channel.",
			// 	Computed:     true,
			// 	NestedObject: ForumTagSchema,
			// },
			"applied_tags": schema.ListAttribute{
				Description: "The IDs of the set of tags that have been applied to a thread in a forum channel.",
				Computed:    true,
				ElementType: types.StringType,
			},
			// "default_reaction_emoji": schema.ObjectAttribute{
			// 	Description: "Emoji to use as the default reaction to a forum post.",
			// 	Computed:    true,
			// 	CustomType:  ForumDefaultReaction,
			// },
			"default_thread_rate_limit_per_user": schema.Int32Attribute{
				Description: "Amount of seconds a user has to wait before sending another message in a thread (0-21600)",
				Computed:    true,
			},
			"default_sort_order": schema.StringAttribute{
				Description: "The default sort order of threads in the channel.",
				Computed:    true,
			},
			"default_forum_layout": schema.StringAttribute{
				Description: "The default layout of threads in the channel.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ChannelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var state, provided ChannelDataSourceModel

	// Read the configuration data into the provided struct.
	diags := req.Config.Get(ctx, &provided)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id": provided.GuildID,
		"name?":    provided.Name,
		"id?":      provided.ID,
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

	var result *discordgo.Channel
	var err error

	// Fetch data from the Discord client
	if id != "" {
		result, err = d.client.Channel(id)
	} else if name != "" {
		result, err = discord.FetchChannelByName(ctx, d.client, guild_id, name)
	} else {
		resp.Diagnostics.AddError(
			"Invalid Resource Configuration",
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

	children, err := discord.FetchChannelChildren(ctx, d.client, guild_id, result)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get children for %s", datasourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	flags := discord.ListStringify(result.Flags)

	flagsList, diags := common.ToListType[string, basetypes.StringType](flags)
	resp.Diagnostics.Append(diags...)

	appliedTags, diags := common.ToListType[string, basetypes.StringType](result.AppliedTags)
	resp.Diagnostics.Append(diags...)

	childrenIDs := discord.ChannelNames(children)
	childrenList, diags := common.ToListType[string, basetypes.StringType](childrenIDs)
	resp.Diagnostics.Append(diags...)

	// availableTags, diags := common.ToListType(channel.AvailableTags)
	// resp.Diagnostics.Append(diags...)

	// attrTypes, attrValues := StructToAttrValues(channel.ThreadMetadata)
	// tflog.Info(ctx, fmt.Sprintf("ThreadMetadata: %v || %v", attrTypes, attrValues))

	// threadMetadata, diags := types.ObjectValueFrom(ctx, attrTypes, channel.ThreadMetadata)
	// resp.Diagnostics.Append(diags...)

	// Map the result data to the state.
	state = ChannelDataSourceModel{
		ID:               types.StringValue(result.ID),
		Type:             types.StringValue(discord.Stringify(result.Type)),
		GuildID:          types.StringValue(result.GuildID),
		Position:         types.Int32Value(int32(result.Position)),
		Name:             types.StringValue(result.Name),
		Topic:            types.StringValue(result.Topic),
		NSFW:             types.BoolValue(result.NSFW),
		Bitrate:          types.Int32Value(int32(result.Bitrate)),
		UserLimit:        types.Int32Value(int32(result.UserLimit)),
		RateLimitPerUser: types.Int32Value(int32(result.RateLimitPerUser)),
		Icon:             types.StringValue(result.Icon),
		OwnerID:          types.StringValue(result.OwnerID),
		ApplicationID:    types.StringValue(result.ApplicationID),
		// Managed:                       types.BoolValue(channel.Managed),
		ParentID:         types.StringValue(result.ParentID),
		Children:         childrenList,
		LastPinTimestamp: types.StringValue(common.StrDiscordTime(result.LastPinTimestamp, "ISO8601")),
		// ThreadMetadata:   threadMetadata,
		// ThreadMember:                  common.ThreadMemberValue(channel.ThreadMember),
		Flags: flagsList,
		// AvailableTags: availableTags,
		AppliedTags: appliedTags,
		// DefaultReactionEmoji:          common.ForumDefaultReactionValue(&channel.DefaultReactionEmoji),
		DefaultThreadRateLimitPerUser: types.Int32Value(int32(result.DefaultThreadRateLimitPerUser)),
		DefaultSortOrder:              types.StringValue(discord.Stringify(result.DefaultSortOrder)),
		DefaultForumLayout:            types.StringValue(discord.Stringify(result.DefaultForumLayout)),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *ChannelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
