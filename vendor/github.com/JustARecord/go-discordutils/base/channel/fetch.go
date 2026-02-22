package channel

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

// FetchByName fetches a channel by name.
func FetchByName(ctx context.Context, client *discordgo.Session, guild_id, name string) (*discordgo.Channel, error) {
	channels, err := client.GuildChannels(guild_id)
	if err != nil {
		return nil, err
	}

	match, ok := lo.Find(channels, func(c *discordgo.Channel) bool {
		return c.Name == name
	})

	if !ok {
		return nil, fmt.Errorf("channel not found: name=%s", name)
	}

	return match, nil
}

// FetchByID fetches a channel by ID.
func FetchByID(ctx context.Context, client *discordgo.Session, guild_id, channel_id string) (*discordgo.Channel, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		channel, err := client.Channel(channel_id)
		if err != nil {
			return nil, err
		}

		return channel, nil
	}
}

// FetchAll fetches all channels in a guild.
func FetchAll(ctx context.Context, client *discordgo.Session, guild_id string) ([]*discordgo.Channel, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.GuildChannels(guild_id)
	}
}

// FetchChildren fetches the children of a channel, if the channel is a category.
func FetchChildren(ctx context.Context, client *discordgo.Session, guild_id string, channel *discordgo.Channel) ([]*discordgo.Channel, error) {
	// If channel type is not a category or a thread, return an empty list
	if channel.Type != discordgo.ChannelTypeGuildCategory && channel.Type != discordgo.ChannelTypeGuildNewsThread && channel.Type != discordgo.ChannelTypeGuildPublicThread && channel.Type != discordgo.ChannelTypeGuildPrivateThread {
		return nil, nil
	}

	channels, err := client.GuildChannels(guild_id)
	if err != nil {
		return nil, err
	}

	children := lo.Filter(channels, func(c *discordgo.Channel, _ int) bool {
		return c.ParentID == channel.ID
	})

	return children, nil
}
