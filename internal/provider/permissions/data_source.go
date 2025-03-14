package permissions

import (
	"context"
	"fmt"

	"github.com/JustARecord/go-discordutils/base/permissions"
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

// NewPermissionsDataSource is a helper function to simplify the provider implementation.
func NewPermissionsDataSource() datasource.DataSource {
	return &PermissionsDataSource{}
}

// Metadata returns the data source type name.
func (d *PermissionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceMetadataName
}

// Schema defines the schema for the data source.
func (d *PermissionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Required:    true,
			},
			"channel_id": schema.StringAttribute{
				Description: "The ID of the channel.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the role or user.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the overwrite, either 'role' or 'member'.",
				Required:    true,
			},
			"allow": schema.ListAttribute{
				Description: "The list of permissions that are allowed.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"deny": schema.ListAttribute{
				Description: "The list of permissions that are denied.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *PermissionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var state, provided PermissionsDataSourceModel

	// Read the configuration data into the provided struct.
	diags := req.Config.Get(ctx, &provided)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id":   provided.GuildID,
		"channel_id": provided.ChannelID,
		"id":         provided.ID,
		"type":       provided.Type,
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
	permissionsType := provided.Type.ValueString()
	guild_id := provided.GuildID.ValueString()
	channel_id := provided.ChannelID.ValueString()

	// Fetch data from the Discord client
	overwrite, err := permissions.FetchChannelPermissions(ctx, d.client, guild_id, channel_id, id, permissionsType)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", datasourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	allowed, denied, err := discord.ParseOverwrite(overwrite)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", datasourceMetadataName),
			err.Error(),
		)
		return
	}

	allowedList, diags := common.ToListType[string, basetypes.StringType](allowed)
	resp.Diagnostics.Append(diags...)

	deniedList, diags := common.ToListType[string, basetypes.StringType](denied)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Map the result data to the state.
	state = PermissionsDataSourceModel{
		GuildID:   types.StringValue(guild_id),
		ChannelID: types.StringValue(channel_id),
		ID:        types.StringValue(id),
		Type:      types.StringValue(permissionsType),
		Allow:     allowedList,
		Deny:      deniedList,
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *PermissionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
