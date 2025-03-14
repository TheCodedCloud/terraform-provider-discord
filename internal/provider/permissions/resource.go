package permissions

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/JustARecord/go-discordutils/base/permissions"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/thecodedcloud/terraform-provider-discord/internal/provider/common"
)

// NewPermissionsResource is a helper function to simplify the provider implementation.
func NewPermissionsResource() resource.Resource {
	return &PermissionsResource{}
}

// Metadata returns the resource type name.
func (r *PermissionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + resourceMetadataName
}

// Schema defines the schema for the resource.
func (r *PermissionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"last_updated": schema.StringAttribute{
				Description: "The last time the resource was updated.",
				Computed:    true,
			},
			"allow": schema.ListAttribute{
				Description: "The list of permissions that are allowed.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"deny": schema.ListAttribute{
				Description: "The list of permissions that are denied.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *PermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx = tflog.SetField(ctx, "operation", "create")

	var plan PermissionsResourceModel

	// Retrieve values from plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id":   plan.GuildID,
		"channel_id": plan.ChannelID,
		"id":         plan.ID,
		"type":       plan.Type,
	}

	// Check for null/unknown values
	common.CheckNonNull(required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	permissionsType := plan.Type.ValueString()
	guild_id := plan.GuildID.ValueString()
	channel_id := plan.ChannelID.ValueString()

	allow, diags := common.FromListType(ctx, plan.Allow)
	resp.Diagnostics.Append(diags...)

	deny, diags := common.FromListType(ctx, plan.Deny)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Sort the allow and deny lists
	sort.Strings(allow)
	sort.Strings(deny)

	// Create the resource
	result, err := permissions.CreatePermissionOverwrite(ctx, r.client, guild_id, channel_id, id, permissionsType, allow, deny)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to create %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading plan %s %s: %v", resourceMetadataName, resourceMetadataType, plan))

	if diags := UpdateModel(ctx, result, &plan, nil); diags != nil {
		resp.Diagnostics.Append(diags...)
	}

	// Set the LastUpdated field to the current time.
	plan.LastUpdated = types.StringValue(common.CurrentTime())

	tflog.Info(ctx, fmt.Sprintf("Updated plan %s %s: %v", resourceMetadataName, resourceMetadataType, plan))

	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *PermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "operation", "update")

	var plan, state PermissionsResourceModel

	// Retrieve values from plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id":   plan.GuildID,
		"channel_id": plan.ChannelID,
		"id":         plan.ID,
		"type":       plan.Type,
	}

	// Check for null/unknown values
	common.CheckNonNull(required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	permissionsType := plan.Type.ValueString()
	guild_id := plan.GuildID.ValueString()
	channel_id := plan.ChannelID.ValueString()

	var allow, deny []string

	if !plan.Allow.IsNull() {
		allow, diags = common.FromListType(ctx, plan.Allow)
		resp.Diagnostics.Append(diags...)
	}

	if !plan.Deny.IsNull() {
		deny, diags = common.FromListType(ctx, plan.Deny)
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Sort the allow and deny lists
	sort.Strings(allow)
	sort.Strings(deny)

	// Update the resource
	result, err := permissions.UpdatePermissionOverwrite(ctx, r.client, guild_id, channel_id, id, permissionsType, allow, deny)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading provided %s %s: %v", resourceMetadataName, resourceMetadataType, plan))
	tflog.Info(ctx, fmt.Sprintf("Reading state %s %s: %v", resourceMetadataName, resourceMetadataType, state))

	if diags := UpdateModel(ctx, result, &plan, nil); diags != nil {
		resp.Diagnostics.Append(diags...)
	}

	// Set the LastUpdated field to the current time.
	plan.LastUpdated = types.StringValue(common.CurrentTime())

	tflog.Info(ctx, fmt.Sprintf("Updated plan %s %s: %v", resourceMetadataName, resourceMetadataType, plan))

	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *PermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx = tflog.SetField(ctx, "operation", "delete")

	var state PermissionsResourceModel

	// Retrieve values from state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id":   state.GuildID,
		"channel_id": state.ChannelID,
		"id":         state.ID,
		"type":       state.Type,
	}

	// Check for null/unknown values
	common.CheckNonNull(required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	permissionsType := state.Type.ValueString()
	guild_id := state.GuildID.ValueString()
	channel_id := state.ChannelID.ValueString()

	// Delete existing resource
	err := permissions.DeletePermissionOverwrite(ctx, r.client, guild_id, channel_id, id, permissionsType)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}
}

// Import imports the resource and sets the Terraform state.
func (r *PermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx = tflog.SetField(ctx, "operation", "import")

	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 4 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <guild_id>/<channel_id>/<type>/<id>. Got: %q", req.ID),
		)
		return
	}

	// First part is the guild ID
	guildID := idParts[0]

	// Second part is the channel ID
	channelID := idParts[1]

	// Third part is the resource type
	resourceType := idParts[2]

	// Fourth part is the resource ID
	resourceID := idParts[3]

	// Set the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), types.StringValue(guildID))...)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("channel_id"), types.StringValue(channelID))...)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("type"), types.StringValue(resourceType))...)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(resourceID))...)
}

// Read refreshes the Terraform state with the latest data.
func (r *PermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var provided PermissionsResourceModel

	// Read the configuration data into the provided struct.
	diags := req.State.Get(ctx, &provided)
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
	common.CheckNonNull(required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	id := provided.ID.ValueString()
	permissionsType := provided.Type.ValueString()
	guild_id := provided.GuildID.ValueString()
	channel_id := provided.ChannelID.ValueString()

	// Fetch data from the Discord client
	result, err := permissions.FetchChannelPermissions(ctx, r.client, guild_id, channel_id, id, permissionsType)
	if err != nil {
		if strings.HasPrefix(err.Error(), "permission overwrite not found") {
			// If the resource is not found, force a recreation and return early
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", datasourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading provided %s %s: %v", resourceMetadataName, resourceMetadataType, provided))

	if diags := UpdateModel(ctx, result, &provided, nil); diags != nil {
		resp.Diagnostics.Append(diags...)
	}

	tflog.Info(ctx, fmt.Sprintf("Updated provided %s %s: %v", resourceMetadataName, resourceMetadataType, provided))

	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	diags = resp.State.Set(ctx, &provided)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *PermissionsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	// Get the client from the provider data.
	client, ok := req.ProviderData.(*discordgo.Session)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unexpected %s Configure Type", resourceMetadataType),
			fmt.Sprintf("Expected *discordgo.Session, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
