package channel

import (
	discord "github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
)

func setupParams(model *ChannelResourceModel) *discordgo.ChannelEdit {
	name := model.Name.ValueString()
	parentID := model.ParentID.ValueString()
	topic := model.Topic.ValueString()

	// Optional parameters
	nsfw := model.NSFW.ValueBool()
	position := int(model.Position.ValueInt32())

	params := &discordgo.ChannelEdit{
		Name:     name,
		ParentID: parentID,
		Topic:    topic,
		NSFW:     &nsfw,
		Position: &position,
	}

	// Set optional parameters
	// TODO: implement

	return params
}

// UpdateModel updates the resource model from the provided data.
func UpdateModel(result *discordgo.Channel, model, state *ChannelResourceModel) diag.Diagnostics {
	// TODO: permissions
	// permissions := common.ListStringifyDiscord(result.Permissions)
	flags := discord.ListStringify(result.Flags)

	// permissionsList, diags := common.ToListType(permissions)
	// if diags.HasError() {
	// 	return diags
	// }

	flagsList, diags := common.ToListType[string, basetypes.StringType](flags)
	if diags.HasError() {
		return diags
	}

	appliedTags, diags := common.ToListType[string, basetypes.StringType](result.AppliedTags)
	if diags.HasError() {
		return diags
	}

	if model == nil {
		model = &ChannelResourceModel{}
	}

	model.ID = types.StringValue(result.ID)
	model.Type = types.StringValue(discord.Stringify(result.Type))
	model.Position = types.Int32Value(int32(result.Position))
	model.Name = types.StringValue(result.Name)
	model.Topic = types.StringValue(result.Topic)
	model.NSFW = types.BoolValue(result.NSFW)
	model.Bitrate = types.Int32Value(int32(result.Bitrate))
	model.UserLimit = types.Int32Value(int32(result.UserLimit))
	model.RateLimitPerUser = types.Int32Value(int32(result.RateLimitPerUser))
	model.Icon = types.StringValue(result.Icon)
	model.OwnerID = types.StringValue(result.OwnerID)
	model.ApplicationID = types.StringValue(result.ApplicationID)
	model.ParentID = types.StringValue(result.ParentID)
	model.LastPinTimestamp = types.StringValue(common.StrDiscordTime(result.LastPinTimestamp, "ISO8601"))
	// model.ThreadMetadata = types.StringValue(result.ThreadMetadata)
	// model.ThreadMember = types.StringValue(result.ThreadMember)
	model.Flags = flagsList
	model.AppliedTags = appliedTags
	// model.DefaultReactionEmoji = types.StringValue(result.DefaultReactionEmoji)
	model.DefaultThreadRateLimitPerUser = types.Int32Value(int32(result.DefaultThreadRateLimitPerUser))
	model.DefaultSortOrder = types.StringValue(discord.Stringify(result.DefaultSortOrder))
	model.DefaultForumLayout = types.StringValue(discord.Stringify(result.DefaultForumLayout))

	if state == nil {
		// If the plan is nil, return early.
		return nil
	}

	// Otherwise, update the model with additional data from the plan.

	// Map the guild data to the state.
	model.GuildID = state.GuildID

	return nil
}
