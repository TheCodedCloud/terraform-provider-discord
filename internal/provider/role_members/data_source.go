package role_members

import (
	"context"
	"fmt"

	"github.com/JustARecord/go-discordutils/base/member"
	"github.com/JustARecord/go-discordutils/base/role"
	"github.com/TheCodedCloud/terraform-provider-discord/internal/provider/common"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// NewRoleMembersDataSource is a helper function to simplify the provider implementation.
func NewRoleMembersDataSource() datasource.DataSource {
	return &RoleMembersDataSource{}
}

// Metadata returns the data source type name.
func (d *RoleMembersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceMetadataName
}

// Schema defines the schema for the data source.
//
// Requires GUILD_MEMBERS: unlike the discord_role_members resource (which was
// changed to check only declared members via the ungated single-member
// endpoint — see internal/provider/discordmembers), this data source has no
// declared-set input to check against. Its only contract is "who in the whole
// guild holds this role," which Discord can only answer via the gated List
// Guild Members endpoint (there is no by-role listing endpoint at any
// privilege level). If the bot application does not have GUILD_MEMBERS
// enabled, Read on this data source will fail with a 403 — that is expected,
// not a bug, and there is currently no ungated way to fix it. See
// docs/PRIVILEGED_INTENTS.md in this repo, and
// docs/data-sources/role_members.md, for the fuller explanation.
func (d *RoleMembersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads the members of a Discord role by enumerating the guild's full member " +
			"list and filtering by role. Requires the GUILD_MEMBERS privileged intent to be " +
			"enabled on the bot application's Developer Portal entry — this is a REST requirement " +
			"of Discord's List Guild Members endpoint, not a gateway/websocket concern, and it " +
			"applies even though this provider never opens a gateway connection. There is no " +
			"Discord endpoint that lists members by role directly, so unlike the " +
			"discord_role_members resource, this data source cannot avoid the gated endpoint by " +
			"checking only declared members.",
		Attributes: map[string]schema.Attribute{
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Required:    true,
			},
			"role_id": schema.StringAttribute{
				Description: "The ID of the role.",
				Optional:    true,
				Computed:    true,
			},
			"role": schema.StringAttribute{
				Description: "The name of the role.",
				Optional:    true,
				Computed:    true,
			},
			"members": schema.ListAttribute{
				Description: "Array of role members. Populating this requires the GUILD_MEMBERS " +
					"privileged intent — see the data source description above.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *RoleMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var state, provided RoleMembersDataSourceModel

	// Read the configuration data into the provided struct.
	diags := req.Config.Get(ctx, &provided)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id": provided.GuildID,
		"role_id?": provided.RoleID,
		"role?":    provided.Role,
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

	role_id := provided.RoleID.ValueString()
	role_name := provided.Role.ValueString()
	guild_id := provided.GuildID.ValueString()

	var result_role *discordgo.Role
	var err error

	// Fetch data from the Discord client
	if role_id != "" {
		result_role, err = role.FetchByID(ctx, d.client, guild_id, role_id)
	} else if role_name != "" {
		result_role, err = role.FetchByName(ctx, d.client, guild_id, role_name)
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s Configuration", datasourceMetadataType),
			fmt.Sprintf("Either the role_id or the role must be set for the %s %s.", datasourceMetadataName, datasourceMetadataType),
		)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", datasourceMetadataName),
			err.Error(),
		)
		return
	}

	// Fetch the role members
	members, err := role.FetchMembers(ctx, d.client, guild_id, result_role.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get roles for %s", datasourceMetadataName),
			err.Error(),
		)
		return
	}

	memberNames := member.Names(members)

	membersList, diags := common.ToListType[string, basetypes.StringType](memberNames)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Map the result data to the state.
	state = RoleMembersDataSourceModel{
		GuildID: types.StringValue(guild_id),
		RoleID:  types.StringValue(result_role.ID),
		Role:    types.StringValue(result_role.Name),
		Members: membersList,
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *RoleMembersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
