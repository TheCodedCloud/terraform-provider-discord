package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// FetchRoleByName fetches a role by name.
func FetchRoleByName(ctx context.Context, client *discordgo.Session, guild, name string) (*discordgo.Role, error) {
	roles, err := client.GuildRoles(guild)
	if err != nil {
		return nil, err
	}

	for _, r := range roles {
		if r.Name == name {
			return r, nil
		}
	}

	return nil, fmt.Errorf("guild not found: %s", name)
}

// FetchRoleByID fetches a role by ID.
func FetchRoleByID(ctx context.Context, client *discordgo.Session, guild, id string) (*discordgo.Role, error) {
	roles, err := client.GuildRoles(guild)
	if err != nil {
		return nil, err
	}

	for _, r := range roles {
		if r.ID == id {
			return r, nil
		}
	}

	return nil, fmt.Errorf("role not found: %s", id)
}
