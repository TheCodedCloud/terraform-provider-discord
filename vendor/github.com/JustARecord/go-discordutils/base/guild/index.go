package guild

import (
	"context"
	"log/slog"
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

// IDs returns a list of guild IDs.
func IDs(ctx context.Context, client *discordgo.Session) ([]string, error) {
	guilds, err := FetchAll(ctx, client)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, g := range guilds {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			ids = append(ids, g.ID)
		}
	}

	// Sort the IDs
	sort.Strings(ids)

	return ids, nil
}

// Names returns a list of guild names.
func Names(ctx context.Context, client *discordgo.Session) ([]string, error) {
	guilds, err := FetchAll(ctx, client)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, g := range guilds {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			names = append(names, g.Name)
		}
	}

	// Sort the names
	sort.Strings(names)

	return names, nil
}

// CanAccess checks if the client can access a guild.
func CanAccess(ctx context.Context, client *discordgo.Session, g *discordgo.Guild) (*bool, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		slog.Info("Checking if the guild is accessible", "guild", g.Name, "id", g.ID)

		// 1. Fetch all guilds that the client has access to.
		guilds, err := FetchAll(ctx, client)
		if err != nil {
			return nil, err
		}

		// 2. Check if the guild is in the list of guilds.
		// a. If the guild ID is not empty, check by ID.
		// b. Otherwise, check by name.
		_, ok := lo.Find(guilds, func(current *discordgo.UserGuild) bool {
			if g.ID != "" {
				return current.ID == g.ID
			}

			return current.Name == g.Name
		})

		return &ok, nil
	}
}
