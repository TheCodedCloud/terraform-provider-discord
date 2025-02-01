package channel

import (
	"context"
	"fmt"
	"strings"

	"github.com/JustARecord/go-discordutils/base/channel"
	discord "github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
)

// NewRoleResource is a helper function to simplify the provider implementation.
func NewRoleResource() resource.Resource {
	return &ChannelResource{}
}

// Metadata returns the resource type name.
func (r *ChannelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + resourceMetadataName
}

// Schema defines the schema for the resource.
func (r *ChannelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"last_updated": schema.StringAttribute{
				Description: "The last time the resource was updated.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the channel.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the channel.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild.",
				Required:    true,
			},
			"position": schema.Int32Attribute{
				Description: "The position of the channel.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the channel.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"topic": schema.StringAttribute{
				Description: "The topic of the channel.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"nsfw": schema.BoolAttribute{
				Description: "Whether the channel is marked as NSFW.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bitrate": schema.Int32Attribute{
				Description: "The bitrate of the channel, if it is a voice channel.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"user_limit": schema.Int32Attribute{
				Description: "The user limit of the voice channel.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"rate_limit_per_user": schema.Int32Attribute{
				Description: "Amount of seconds a user has to wait before sending another message or creating another thread (0-21600)",
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"icon": schema.StringAttribute{
				Description: "Icon of the group DM channel.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner_id": schema.StringAttribute{
				Description: "ID of the creator of the group DM or thread",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				Description: "ApplicationID of the DM creator Zeroed if guild channel or not a bot user",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// "managed": schema.BoolAttribute{
			// 	Description: "Whether the role is managed by an integration.",
			// 	Computed:    true,
			// },
			"parent_id": schema.StringAttribute{
				Description: "The ID of the parent category for a channel.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"children": schema.ListAttribute{
				Description: "The IDs of the child channels of the category, if the channel is a category.",
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"last_pin_timestamp": schema.StringAttribute{
				Description: "The timestamp of the last pinned message in the channel.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			// "default_reaction_emoji": schema.ObjectAttribute{
			// 	Description: "Emoji to use as the default reaction to a forum post.",
			// 	Computed:    true,
			// 	CustomType:  ForumDefaultReaction,
			// },
			"default_thread_rate_limit_per_user": schema.Int32Attribute{
				Description: "Amount of seconds a user has to wait before sending another message in a thread (0-21600)",
				Computed:    true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"default_sort_order": schema.StringAttribute{
				Description: "The default sort order of threads in the channel.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_forum_layout": schema.StringAttribute{
				Description: "The default layout of threads in the channel.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx = tflog.SetField(ctx, "operation", "create")

	var plan ChannelResourceModel

	// Retrieve values from plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id": plan.GuildID,
		"name":     plan.Name,
	}

	// Check for null/unknown values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	guild_id := plan.GuildID.ValueString()
	name := plan.Name.ValueString()
	channelTypeStr := plan.Type.ValueString()

	params := setupParams(&plan)

	// Create the resource
	result, err := channel.CreateWithParams(ctx, r.client, guild_id, name, channelTypeStr, params)
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

	if diags := UpdateModel(result, &plan, nil); diags != nil {
		resp.Diagnostics.Append(diags...)
	}

	children, err := channel.FetchChildren(ctx, r.client, plan.GuildID.ValueString(), result)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get children for %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	childrenIDs := channel.Names(children)
	childrenList, diags := common.ToListType[string, basetypes.StringType](childrenIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the children
	plan.Children = childrenList

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
func (r *ChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "operation", "update")

	var plan, state ChannelResourceModel

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
		"name?":    plan.Name,
		"id?":      plan.ID,
	}

	// Check for null/unknown values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()

	params := setupParams(&plan)

	// Update the resource
	result, err := channel.UpdateByID(ctx, r.client, id, params)
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

	if diags := UpdateModel(result, &plan, nil); diags != nil {
		resp.Diagnostics.Append(diags...)
	}

	children, err := channel.FetchChildren(ctx, r.client, plan.GuildID.ValueString(), result)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get children for %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	childrenIDs := channel.Names(children)
	childrenList, diags := common.ToListType[string, basetypes.StringType](childrenIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the children
	plan.Children = childrenList

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
func (r *ChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx = tflog.SetField(ctx, "operation", "delete")

	var state ChannelResourceModel

	// Retrieve values from state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"guild_id": state.GuildID,
		"name?":    state.Name,
		"id?":      state.ID,
	}

	// Check for null/unknown values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required values
	common.CheckRequired(ctx, required, resp.Diagnostics, resourceMetadataName, resourceMetadataType)
	if resp.Diagnostics.HasError() {
		return
	}

	guild_id := state.GuildID.ValueString()

	var err error

	// Delete existing resource
	if !state.ID.IsNull() {
		err = channel.DeleteByID(ctx, r.client, state.ID.ValueString())
	} else if !state.Name.IsNull() {
		err = channel.DeleteByName(ctx, r.client, guild_id, state.Name.ValueString())
	} else {
		err = fmt.Errorf("either the id or the name must be set for the %s %s", resourceMetadataName, resourceMetadataType)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete %s", resourceMetadataName),
			err.Error(),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
}

// Import imports the resource and sets the Terraform state.
func (r *ChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

	// Second part is the id part
	resourcePart := idParts[1]

	// Set the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), types.StringValue(guildID))...)

	// Check if the id part is an ID or a name
	// If ID is a snowflake, it's an ID
	if discord.IsSnowflake(resourcePart) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(resourcePart))...)
	} else {
		// Otherwise, it's a name
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), types.StringValue(resourcePart))...)
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *ChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var state, provided ChannelResourceModel

	// Read the configuration data into the provided struct.
	diags := req.State.Get(ctx, &provided)
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
	name := provided.Name.ValueString()
	guild_id := provided.GuildID.ValueString()

	var result *discordgo.Channel
	var err error

	// Fetch data from the Discord client
	if id != "" {
		result, err = channel.FetchByID(ctx, r.client, guild_id, id)
	} else if name != "" {
		result, err = channel.FetchByName(ctx, r.client, guild_id, name)
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s Configuration", resourceMetadataType),
			fmt.Sprintf("Either the id or the name must be set for the %s %s.", resourceMetadataName, resourceMetadataType),
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

	tflog.Info(ctx, fmt.Sprintf("Reading provided %s %s: %v", resourceMetadataName, resourceMetadataType, provided))
	tflog.Info(ctx, fmt.Sprintf("Reading state %s %s: %v", resourceMetadataName, resourceMetadataType, state))

	if diags := UpdateModel(result, &state, &provided); diags != nil {
		resp.Diagnostics.Append(diags...)
	}

	tflog.Info(ctx, fmt.Sprintf("Updated state %s %s: %v", resourceMetadataName, resourceMetadataType, state))

	if resp.Diagnostics.HasError() {
		return
	}

	children, err := channel.FetchChildren(ctx, r.client, state.GuildID.ValueString(), result)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get children for %s", resourceMetadataName),
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	childrenIDs := channel.Names(children)
	childrenList, diags := common.ToListType[string, basetypes.StringType](childrenIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the children
	state.Children = childrenList

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
func (r *ChannelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
