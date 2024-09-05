package channel

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
	"github.com/justarecord/terraform-provider-discord/internal/provider/discord"
)

func setupParams(model *ChannelResourceModel) *discordgo.ChannelEdit {
	name := model.Name.ValueString()
	parentID := model.ParentID.ValueString()
	topic := model.Topic.ValueString()

	params := &discordgo.ChannelEdit{
		Name:     name,
		ParentID: parentID,
		Topic:    topic,
	}

	// Set optional parameters
	// TODO: implement

	return params
}

// CreateChannel creates a new channel with the provided model
func CreateChannel(ctx context.Context, client *discordgo.Session, model *ChannelResourceModel) (*discordgo.Channel, error) {
	guild_id := model.GuildID.ValueString()
	name := model.Name.ValueString()
	channelTypeStr := model.Type.ValueString()

	if channelTypeStr == "" {
		channelTypeStr = discord.Stringify(discordgo.ChannelTypeGuildText)
	}

	channelType, ok := common.FetchKeyByValue(discord.ChannelTypes, channelTypeStr)
	if !ok {
		return nil, fmt.Errorf("invalid channel type %s", channelTypeStr)
	}

	// Setup the parameters
	// TODO: implement
	// params := setupParams(model, permissions)

	// Create it
	result, err := client.GuildChannelCreate(guild_id, name, channelType.(discordgo.ChannelType))
	if err != nil {
		return nil, err
	}

	// We also need to update the channel with the rest of the data.
	// This is because we can't set all the data in the initial creation.

	params := setupParams(model)

	result, err = client.ChannelEdit(result.ID, params)
	if err != nil {
		return result, err
	}

	return result, nil
}

// UpdateChannel updates the channel with the provided model.
func UpdateChannel(ctx context.Context, client *discordgo.Session, model *ChannelResourceModel) (*discordgo.Channel, error) {
	id := model.ID.ValueString()

	// Setup the parameters
	params := setupParams(model)

	// Update it
	result, err := client.ChannelEdit(id, params)
	if err != nil {
		return nil, err
	}

	return result, nil
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
	model.LastMessageID = types.StringValue(result.LastMessageID)
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

// DeleteChannel deletes the channel with the provided state.
func DeleteChannel(ctx context.Context, client *discordgo.Session, model *ChannelResourceModel) error {
	var err error

	guild_id := model.GuildID.ValueString()

	// If the ID is set, delete the role by ID.
	if !model.ID.IsNull() {
		_, err = client.ChannelDelete(model.ID.ValueString())
	} else if !model.Name.IsNull() {
		var channel *discordgo.Channel
		// If the ID is not set, delete the channel by name.

		// Likely, we shouldn't reach this point as the channel wouldn't be in the Terraform state,
		// but it's here for completeness.
		channel, err = discord.FetchChannelByName(ctx, client, guild_id, model.Name.ValueString())
		if err != nil {
			return err
		}

		_, err = client.ChannelDelete(channel.ID)
	}

	return err
}
