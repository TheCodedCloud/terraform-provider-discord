package role

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

// NewRoleDataSource is a helper function to simplify the provider implementation.
func NewRoleDataSource() datasource.DataSource {
	return &RoleDataSource{}
}

// Metadata returns the data source type name.
func (d *RoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceMetadataName
}

// Schema defines the schema for the data source.
func (d *RoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the role.",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the role.",
				Optional:    true,
				Computed:    true,
			},
			"managed": schema.BoolAttribute{
				Description: "Whether this role is managed by an integration, and thus cannot be manually added to, or taken from, members.",
				Computed:    true,
			},
			"mentionable": schema.BoolAttribute{
				Description: "Whether this role is mentionable.",
				Computed:    true,
			},
			"hoist": schema.BoolAttribute{
				Description: "Whether this role is hoisted (shows up separately in member list).",
				Computed:    true,
			},
			"color": schema.StringAttribute{
				Description: "The hex color of this role.",
				Computed:    true,
			},
			"position": schema.Int32Attribute{
				Description: "The position of this role in the guild's role hierarchy.",
				Computed:    true,
			},
			"permissions": schema.ListAttribute{
				Description: "The permissions of the role on the guild (doesn't include channel overrides).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"icon": schema.StringAttribute{
				Description: "The hash of the role icon. Use Role.IconURL to retrieve the icon's URL.",
				Computed:    true,
			},
			"unicode_emoji": schema.StringAttribute{
				Description: "The emoji assigned to this role.",
				Computed:    true,
			},
			"flags": schema.ListAttribute{
				Description: "The flags of the role, which describe its extra features.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *RoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var state, provided RoleDataSourceModel

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

	var result *discordgo.Role
	var err error

	// Fetch data from the Discord client
	if id != "" {
		result, err = discord.FetchRoleByID(ctx, d.client, guild_id, id)
	} else if name != "" {
		result, err = discord.FetchRoleByName(ctx, d.client, guild_id, name)
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
		return
	}

	permissions := discord.ListStringify(result.Permissions)
	flags := discord.ListStringify(result.Flags)

	permissionsList, diags := common.ToListType[string, basetypes.StringType](permissions)
	resp.Diagnostics.Append(diags...)

	flagsList, diags := common.ToListType[string, basetypes.StringType](flags)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Map the result data to the state.
	state = RoleDataSourceModel{
		GuildID:      types.StringValue(guild_id),
		ID:           types.StringValue(result.ID),
		Name:         types.StringValue(result.Name),
		Managed:      types.BoolValue(result.Managed),
		Mentionable:  types.BoolValue(result.Mentionable),
		Hoist:        types.BoolValue(result.Hoist),
		Color:        types.StringValue(common.StrHex(result.Color)),
		Position:     types.Int32Value(int32(result.Position)),
		Permissions:  permissionsList,
		Icon:         types.StringValue(result.Icon),
		UnicodeEmoji: types.StringValue(result.UnicodeEmoji),
		Flags:        flagsList,
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *RoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
