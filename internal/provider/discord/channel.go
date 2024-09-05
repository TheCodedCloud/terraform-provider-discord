package discord

import (
	"context"
	"fmt"

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
