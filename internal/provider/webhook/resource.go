package webhook

import (
	"context"
	"fmt"
	"strings"

	"github.com/JustARecord/go-discordutils/base/channel"
	"github.com/JustARecord/go-discordutils/base/webhook"
	discord "github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
)

// NewWebhookResource is a helper function to simplify the provider implementation.
func NewWebhookResource() resource.Resource {
	return &WebhookResource{}
}

// Metadata returns the resource type name.
func (r *WebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + resourceMetadataName
}

// Schema defines the schema for the resource.
func (r *WebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"last_updated": schema.StringAttribute{
				Description: "The last time the resource was updated.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the webhook.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the webhook, either 'INCOMING', 'CHANNEL_FOLLOWER', or 'APPLICATION'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guild_id": schema.StringAttribute{
				Description: "The guild ID this webhook is for, if any.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel_id": schema.StringAttribute{
				Description: "The channel ID this webhook is for, if any.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The default name of the webhook.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"avatar": schema.StringAttribute{
				Description: "The default user avatar hash of the webhook.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true, // Marked as sensitive to prevent log spam with large base64 encoded images.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token": schema.StringAttribute{
				Description: "The secure token of the webhook (returned for Incoming Webhooks).",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				Description: "The Bot/OAuth2 application that created this webhook.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *WebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx = tflog.SetField(ctx, "operation", "create")

	var plan WebhookResourceModel

	// Retrieve values from plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"channel_id": plan.ChannelID,
		"name":       plan.Name,
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

	channel_id := plan.ChannelID.ValueString()
	name := plan.Name.ValueString()

	// Optional fields, default to "" if not set
	avatar := plan.Avatar.ValueString()

	// Create the resource
	result, err := webhook.CreateWebhook(ctx, r.client, channel_id, name, avatar)
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
func (r *WebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "operation", "update")

	var plan, state WebhookResourceModel

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
		"id": plan.ID,
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

	// The rest of the fields are optional and default to "" if not set.
	name := plan.Name.ValueString()
	avatar := plan.Avatar.ValueString()
	channelId := plan.ChannelID.ValueString()

	// Update the resource
	result, err := webhook.UpdateWebhook(ctx, r.client, id, name, avatar, channelId)
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
func (r *WebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx = tflog.SetField(ctx, "operation", "delete")

	var state WebhookResourceModel

	// Retrieve values from state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	required := map[string]attr.Value{
		"id": state.ID,
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

	// Delete existing resource
	err := webhook.DeleteWebhook(ctx, r.client, id)
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
func (r *WebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx = tflog.SetField(ctx, "operation", "import")

	idParts := strings.Split(req.ID, "/")

	// There are two possible import identifiers:
	// 1. <channel|guild>/<channel_id|guild_id>/<name>
	// 2. <id>

	var guildID, channelID, importType, resource string

	// Check if the import identifier is a single ID
	if len(idParts) == 1 {
		resource = idParts[0]
		importType = "id"
	} else {
		if idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier",
				fmt.Sprintf("Expected import identifier with format: <channel|guild>/<channel_id|guild_id>/<name> or <id>. Got: %q", req.ID),
			)
			return
		}

		// First part is the importType (guild or channel)
		importType = idParts[0]
		if importType != "guild" && importType != "channel" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier",
				fmt.Sprintf("Expected import identifier with format: <channel|guild>/<channel_id|guild_id>/<name> or <id>. Got: %q", req.ID),
			)
			return
		}

		// Second part is the guild ID or channel ID
		if importType == "guild" {
			guildID = idParts[1]
		} else {
			channelID = idParts[1]
		}

		// Third part is the name
		resource = idParts[2]
	}

	// Set the state
	if importType == "guild" {
		// Set the guild ID and resource name
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), types.StringValue(guildID))...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), types.StringValue(resource))...)
	} else if importType == "channel" {
		if discord.IsSnowflake(channelID) {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("channel_id"), types.StringValue(channelID))...)
		} else {

			guildID = idParts[1]
			channelID = idParts[2]
			resource = idParts[3]

			// Input is a channel name, fetch the channel ID
			channel, err := channel.FetchByName(ctx, r.client, guildID, channelID)
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to import state",
					err.Error(),
				)
				return
			}

			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("channel_id"), types.StringValue(channel.ID))...)
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), types.StringValue(resource))...)
	} else if importType == "id" {
		// Set the resource ID
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(resource))...)
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *WebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "operation", "read")

	var provided WebhookResourceModel

	// Read the configuration data into the provided struct.
	diags := req.State.Get(ctx, &provided)
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
	channel_id := provided.ChannelID.ValueString()

	var result *discordgo.Webhook
	var err error

	// Fetch data from the Discord client
	if id != "" {
		result, err = webhook.FetchByID(ctx, r.client, id)
	} else if name != "" {
		if channel_id != "" {
			result, err = webhook.FetchChannelWebhookByName(ctx, r.client, channel_id, name)
		} else if guild_id != "" {
			result, err = webhook.FetchGuildWebhookByName(ctx, r.client, guild_id, name)
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
func (r *WebhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
