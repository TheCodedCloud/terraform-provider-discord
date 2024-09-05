package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// FetchGuildByName fetches a guild by name.
func FetchGuildByName(ctx context.Context, client *discordgo.Session, name string) (*discordgo.Guild, error) {
	guilds, err := client.UserGuilds(100, "", "", false)
	if err != nil {
		return nil, err
	}

	for _, g := range guilds {
		if g.Name == name {
			return FetchGuildByID(ctx, client, g.ID)
		}
	}

	return nil, fmt.Errorf("guild not found: %s", name)
}

// FetchGuildByID fetches a guild by ID.
func FetchGuildByID(ctx context.Context, client *discordgo.Session, id string) (*discordgo.Guild, error) {
	return client.Guild(id)
}
