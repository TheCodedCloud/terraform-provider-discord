package discord

import (
	"context"
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
)

// FetchChannelByName fetches a channel by name.
func FetchChannelByName(ctx context.Context, client *discordgo.Session, guild_id, name string) (*discordgo.Channel, error) {
	channels, err := client.GuildChannels(guild_id)
	if err != nil {
		return nil, err
	}

	for _, c := range channels {
		if c.Name == name {
			return c, nil
		}
	}

	return nil, fmt.Errorf("channel not found: %s", name)
}

// FetchChannelByID fetches a channel by ID.
func FetchChannelByID(ctx context.Context, client *discordgo.Session, guild_id, channel_id string) (*discordgo.Channel, error) {
	channel, err := client.Channel(channel_id)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

// FetchChannelChildren fetches the children of a channel, if the channel is a category.
func FetchChannelChildren(ctx context.Context, client *discordgo.Session, guild_id string, channel *discordgo.Channel) ([]*discordgo.Channel, error) {
	var children []*discordgo.Channel

	// If channel type is not a category or a thread, return an empty list
	if channel.Type != discordgo.ChannelTypeGuildCategory && channel.Type != discordgo.ChannelTypeGuildNewsThread && channel.Type != discordgo.ChannelTypeGuildPublicThread && channel.Type != discordgo.ChannelTypeGuildPrivateThread {
		return children, nil
	}

	channels, err := client.GuildChannels(guild_id)
	if err != nil {
		return nil, err
	}

	for _, c := range channels {
		if c.ParentID == channel.ID {
			children = append(children, c)
		}
	}

	return children, nil
}

// ChannelIDs returns a list of channel IDs.
func ChannelIDs(channels []*discordgo.Channel) []string {
	ids := make([]string, len(channels))
	for i, c := range channels {
		ids[i] = c.ID
	}

	// Sort the IDs
	sort.Strings(ids)

	return ids
}

// ChannelNames returns a list of channel names.
func ChannelNames(channels []*discordgo.Channel) []string {
	names := make([]string, len(channels))
	for i, c := range channels {
		names[i] = c.Name
	}

	// Sort the names
	sort.Strings(names)

	return names
}
