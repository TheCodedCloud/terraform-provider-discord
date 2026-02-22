package channel

import (
	"context"
	"fmt"

	"github.com/JustARecord/go-discordutils/base/common"
	"github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
)

// Create creates a channel.
func Create(ctx context.Context, client *discordgo.Session, guild_id, name, channel_type_str string) (*discordgo.Channel, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if channel_type_str == "" {
			channel_type_str = utils.Stringify(discordgo.ChannelTypeGuildText)
		}

		channel_type, ok := utils.KeyStringify[discordgo.ChannelType](common.ChannelType, channel_type_str)
		if !ok {
			return nil, fmt.Errorf("invalid channel type: %s", channel_type_str)
		}

		return client.GuildChannelCreate(guild_id, name, channel_type.(discordgo.ChannelType))
	}
}

// Update updates a channel.
func Update(ctx context.Context, client *discordgo.Session, channel *discordgo.Channel, params *discordgo.ChannelEdit) (*discordgo.Channel, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.ChannelEdit(channel.ID, params)
	}
}

// UpdateByID updates a channel by ID.
func UpdateByID(ctx context.Context, client *discordgo.Session, channel_id string, params *discordgo.ChannelEdit) (*discordgo.Channel, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.ChannelEdit(channel_id, params)
	}
}

// UpdateByName updates a channel by name.
func UpdateByName(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, name string, params *discordgo.ChannelEdit) (*discordgo.Channel, error) {
	channel, err := FetchByName(ctx, client, guild.ID, name)
	if err != nil {
		return nil, err
	}

	return Update(ctx, client, channel, params)
}

// Delete deletes a channel.
func Delete(ctx context.Context, client *discordgo.Session, channel *discordgo.Channel) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := client.ChannelDelete(channel.ID)
		return err
	}
}

// DeleteByID deletes a channel by ID.
func DeleteByID(ctx context.Context, client *discordgo.Session, channel_id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := client.ChannelDelete(channel_id)
		return err
	}
}

// DeleteByName deletes a channel by name.
func DeleteByName(ctx context.Context, client *discordgo.Session, guild_id, name string) error {
	channel, err := FetchByName(ctx, client, guild_id, name)
	if err != nil {
		return err
	}

	return Delete(ctx, client, channel)
}

// CreateWithParams creates a channel with parameters.
func CreateWithParams(ctx context.Context, client *discordgo.Session, guild_id, name, channel_type string, params *discordgo.ChannelEdit) (*discordgo.Channel, error) {
	// 1. Create a channel.
	channel, err := Create(ctx, client, guild_id, name, channel_type)
	if err != nil {
		return nil, err
	}

	// 2. Update the channel.
	// TODO: replace this all with just CreateComplex
	return Update(ctx, client, channel, params)
}
