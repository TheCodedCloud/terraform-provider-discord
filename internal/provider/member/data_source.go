package member

import (
	"context"
	"fmt"

	"github.com/JustARecord/go-discordutils/base/member"
	"github.com/JustARecord/go-discordutils/base/role"
	discord "github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/thecodedcloud/terraform-provider-discord/internal/provider/common"
)

// NewMemberDataSource is a helper function to simplify the provider implementation.
func NewMemberDataSource() datasource.DataSource {
	return &MemberDataSource{}
}

// Metadata returns the data source type name.
func (d *MemberDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceMetadataName
}

// Schema defines the schema for the data source.
func (d *MemberDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the member.",
				Optional:    true,
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username of the member, not unique across the platform.",
				Optional:    true,
				Computed:    true,
			},
			"discriminator": schema.StringAttribute{
				Description: "The user's Discord-tag.",
				Optional:    true,
				Computed:    true,
			},
			"nick": schema.StringAttribute{
				Description: "This user's guild nickname (if one is set).",
				Computed:    true,
			},
			"avatar": schema.StringAttribute{
				Description: "The member's guild avatar hash.",
				Computed:    true,
			},
			"user": UserSchema,
			"roles": schema.ListAttribute{
				Description: "Array of roles",
				Computed:    true,
				ElementType: types.StringType,
			},
			"joined_at": schema.StringAttribute{
				Description: "When the user joined the guild.",
				Computed:    true,
			},
			"premium_since": schema.StringAttribute{
				Description: "When the user started boosting the guild.",
				Computed:    true,
			},
			"deaf": schema.BoolAttribute{
				Description: "Whether the user is deafened in voice channels.",
				Computed:    true,
			},
			"mute": schema.BoolAttribute{
				Description: "Whether the user is muted in voice channels.",
				Computed:    true,
			},
			"flags": schema.ListAttribute{
				Description: "Guild member flags",
				Computed:    true,
				ElementType: types.StringType,
			},
			"pending": schema.BoolAttribute{
				Description: "Whether the user has not yet passed the guild's Membership Screening requirements.",
				Computed:    true,
			},
			"permissions": schema.ListAttribute{
				Description: "Total permissions of the member in the channel, including overwrites, returned when in the interaction object.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"communication_disabled_until": schema.StringAttribute{
				Description: "When the user's timeout will expire.",
				Computed:    true,
			},
			// "avatar_decoration_data": schema.StringAttribute{
			// 	Description: "Data for the member's guild avatar decoration.",
			// 	Computed:    true,
			// },
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *MemberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var state, provided MemberDataSourceModel

	// Read the configuration data into the provided struct.
	diags := req.Config.Get(ctx, &provided)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id":  provided.GuildID,
		"username?": provided.Username,
		"id?":       provided.ID,
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
	username := provided.Username.ValueString()
	guild_id := provided.GuildID.ValueString()

	var result *discordgo.Member
	var err error

	// Fetch data from the Discord client
	if id != "" {
		result, err = member.FetchById(ctx, d.client, guild_id, id)
	} else if username != "" {
		result, err = member.FetchByName(ctx, d.client, guild_id, username)
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s Configuration", datasourceMetadataType),
			fmt.Sprintf("Either the id or the username must be set for the %s %s.", datasourceMetadataName, datasourceMetadataType),
		)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", datasourceMetadataName),
			err.Error(),
		)
		return
	}

	// Fetch the roles from Discord
	roles, err := role.FetchByIDs(ctx, d.client, guild_id, result.Roles)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get roles for %s", datasourceMetadataName),
			err.Error(),
		)
		return
	}

	// Sort the roles by Discord role position
	role.SortByPosition(roles)

	roleNames := role.Names(roles)

	flags := discord.ListStringify(result.Flags)
	permissions := discord.ListStringify(result.Permissions)

	rolesList, diags := common.ToListType[string, basetypes.StringType](roleNames)
	resp.Diagnostics.Append(diags...)

	flagsList, diags := common.ToListType[string, basetypes.StringType](flags)
	resp.Diagnostics.Append(diags...)

	permissionsList, diags := common.ToListType[string, basetypes.StringType](permissions)
	resp.Diagnostics.Append(diags...)

	user, diags := ToUser(ctx, result.User, guild_id)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Map the result data to the state.
	state = MemberDataSourceModel{
		GuildID:                    types.StringValue(guild_id),
		ID:                         types.StringValue(result.User.ID),
		Username:                   types.StringValue(result.User.Username),
		Discriminator:              types.StringValue(result.User.Discriminator),
		Nick:                       types.StringValue(result.Nick),
		Avatar:                     types.StringValue(result.Avatar),
		User:                       user,
		Roles:                      rolesList,
		JoinedAt:                   types.StringValue(common.StrDiscordTime(&result.JoinedAt, "ISO8601")),
		PremiumSince:               types.StringValue(common.StrDiscordTime(result.PremiumSince, "ISO8601")),
		Deaf:                       types.BoolValue(result.Deaf),
		Mute:                       types.BoolValue(result.Mute),
		Flags:                      flagsList,
		Pending:                    types.BoolValue(result.Pending),
		Permissions:                permissionsList,
		CommunicationDisabledUntil: types.StringValue(common.StrDiscordTime(result.CommunicationDisabledUntil, "ISO8601")),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *MemberDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
