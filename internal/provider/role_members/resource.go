package role_members

import (
	"context"
	"fmt"
	"strings"

	"github.com/JustARecord/go-discordutils/base/guild"
	"github.com/JustARecord/go-discordutils/base/role"
	discord "github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
)

// NewRoleMembersResource is a helper function to simplify the provider implementation.
func NewRoleMembersResource() resource.Resource {
	return &RoleMembersResource{}
}

// Metadata returns the resource type name.
func (r *RoleMembersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + resourceMetadataName
}

// Schema defines the schema for the resource.
func (r *RoleMembersResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_id": schema.StringAttribute{
				Description: "The ID of the role.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role": schema.StringAttribute{
				Description: "The name of the role.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "The last time the resource was updated.",
				Computed:    true,
			},
			"members": schema.ListAttribute{
				Description: "Array of role members",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *RoleMembersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx = tflog.SetField(ctx, "operation", "create")

	var plan RoleMembersResourceModel

	// Retrieve values from plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id": plan.GuildID,
		"role_id?": plan.RoleID,
		"role?":    plan.Role,
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

	members_names := []string{}

	// Set the members
	if !plan.Members.IsNull() {
		members_names, diags = common.FromListType(ctx, plan.Members)

		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	guild_id := plan.GuildID.ValueString()

	// Fetch the guild
	result_guild, err := guild.FetchByID(ctx, r.client, guild_id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the role
	var result_role *discordgo.Role

	if !plan.RoleID.IsNull() {
		result_role, err = role.FetchByID(ctx, r.client, guild_id, plan.RoleID.ValueString())
	} else if !plan.Role.IsNull() {
		result_role, err = role.FetchByName(ctx, r.client, guild_id, plan.Role.ValueString())
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s Configuration", resourceMetadataType),
			fmt.Sprintf("Either the role_id or the role must be set for the %s %s.", resourceMetadataName, resourceMetadataType),
		)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the members
	members, err := guild.FetchMembersByName(ctx, r.client, result_guild, members_names)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get members for %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the resource
	result, err := role.SetMembers(ctx, r.client, result_guild, result_role, members)
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

	if diags := UpdateModel(result_role, result, &plan, nil); diags != nil {
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
func (r *RoleMembersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "operation", "update")

	var plan, state RoleMembersResourceModel

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
		"guild_id": plan.GuildID,
		"role_id?": plan.RoleID,
		"role?":    state.Role,
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

	members_names := []string{}

	// Set the members
	if !plan.Members.IsNull() {
		members_names, diags = common.FromListType(ctx, plan.Members)

		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	guild_id := plan.GuildID.ValueString()

	// Fetch the guild
	result_guild, err := guild.FetchByID(ctx, r.client, guild_id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the role
	var result_role *discordgo.Role

	if !plan.RoleID.IsNull() {
		result_role, err = role.FetchByID(ctx, r.client, guild_id, plan.RoleID.ValueString())
	} else if !state.Role.IsNull() {
		result_role, err = role.FetchByName(ctx, r.client, guild_id, state.Role.ValueString())
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s Configuration", resourceMetadataType),
			fmt.Sprintf("Either the role_id or the role must be set for the %s %s.", resourceMetadataName, resourceMetadataType),
		)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the members
	members, err := guild.FetchMembersByName(ctx, r.client, result_guild, members_names)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get members for %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the resource
	result, err := role.SetMembers(ctx, r.client, result_guild, result_role, members)
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

	if diags := UpdateModel(result_role, result, &plan, nil); diags != nil {
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
func (r *RoleMembersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx = tflog.SetField(ctx, "operation", "delete")

	var state RoleMembersResourceModel

	// Retrieve values from state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id": state.GuildID,
		"role_id?": state.RoleID,
		"role?":    state.Role,
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

	members_names := []string{}

	// Set the members
	if !state.Members.IsNull() {
		members_names, diags = common.FromListType(ctx, state.Members)

		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	guild_id := state.GuildID.ValueString()

	// Fetch the guild
	result_guild, err := guild.FetchByID(ctx, r.client, guild_id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the role
	var result_role *discordgo.Role

	if !state.RoleID.IsNull() {
		result_role, err = role.FetchByID(ctx, r.client, guild_id, state.RoleID.ValueString())
	} else if !state.Role.IsNull() {
		result_role, err = role.FetchByName(ctx, r.client, guild_id, state.Role.ValueString())
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s Configuration", resourceMetadataType),
			fmt.Sprintf("Either the role_id or the role must be set for the %s %s.", resourceMetadataName, resourceMetadataType),
		)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the members
	members, err := guild.FetchMembersByName(ctx, r.client, result_guild, members_names)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get members for %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove the members
	if err = role.RemoveMembers(ctx, r.client, result_guild, result_role, members); err != nil {
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
func (r *RoleMembersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx = tflog.SetField(ctx, "operation", "import")

	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <guild_id>/<id|name>. Got: %q", req.ID),
		)
		return
	}

	// First part is the guild ID
	guildID := idParts[0]

	// Second part is the resource part
	resourcePart := idParts[1]

	// Set the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), types.StringValue(guildID))...)

	// Check if the role part is an ID or a name
	// If ID is a snowflake, it's an ID
	if discord.IsSnowflake(resourcePart) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(resourcePart))...)
	} else {
		// Otherwise, it's a name
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), types.StringValue(resourcePart))...)
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *RoleMembersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var state, provided RoleMembersResourceModel

	// Read the configuration data into the provided struct.
	diags := req.State.Get(ctx, &provided)
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
	common.CheckNonNull(required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	guild_id := provided.GuildID.ValueString()

	// Fetch the guild
	result_guild, err := guild.FetchByID(ctx, r.client, guild_id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the role
	var result_role *discordgo.Role
	if !provided.RoleID.IsNull() {
		result_role, err = role.FetchByID(ctx, r.client, guild_id, provided.RoleID.ValueString())
	} else if !provided.Role.IsNull() {
		result_role, err = role.FetchByName(ctx, r.client, guild_id, provided.Role.ValueString())
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s Configuration", resourceMetadataType),
			fmt.Sprintf("Either the role_id or the role must be set for the %s %s.", resourceMetadataName, resourceMetadataType),
		)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Read the resource
	result, err := role.FetchMembers(ctx, r.client, result_guild.ID, result_role.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read %s", resourceMetadataName),
			err.Error(),
		)
	}

	tflog.Info(ctx, fmt.Sprintf("Reading provided %s %s: %v", resourceMetadataName, resourceMetadataType, provided))
	tflog.Info(ctx, fmt.Sprintf("Reading state %s %s: %v", resourceMetadataName, resourceMetadataType, state))

	if diags := UpdateModel(result_role, result, &state, &provided); diags != nil {
		resp.Diagnostics.Append(diags...)
	}

	tflog.Info(ctx, fmt.Sprintf("Updated state %s %s: %v", resourceMetadataName, resourceMetadataType, state))

	if resp.Diagnostics.HasError() {
		return
	}

	// Revert last_updated to the plan value
	if !provided.LastUpdated.IsNull() {
		state.LastUpdated = provided.LastUpdated
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *RoleMembersResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
