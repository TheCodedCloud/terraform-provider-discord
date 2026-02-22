package guild

import (
	"context"
	"fmt"

	"github.com/JustARecord/go-discordutils/base/member"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

// FetchAll fetches all guilds.
func FetchAll(ctx context.Context, client *discordgo.Session) ([]*discordgo.UserGuild, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.UserGuilds(100, "", "", false)
	}
}

// FetchByName fetches a guild by name.
func FetchByName(ctx context.Context, client *discordgo.Session, name string) (*discordgo.Guild, error) {
	guilds, err := FetchAll(ctx, client)
	if err != nil {
		return nil, err
	}

	match, ok := lo.Find(guilds, func(g *discordgo.UserGuild) bool {
		return g.Name == name
	})

	if !ok {
		return nil, fmt.Errorf("guild not found: name=%s", name)
	}

	return FetchByID(ctx, client, match.ID)
}

// FetchByID fetches a guild by ID.
func FetchByID(ctx context.Context, client *discordgo.Session, id string) (*discordgo.Guild, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.Guild(id)
	}
}

// FetchMembersByName fetches members of a guild by name.
func FetchMembersByName(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, names []string) ([]*discordgo.Member, error) {
	members, err := member.All(ctx, client, guild.ID)
	if err != nil {
		return nil, err
	}

	result := lo.Filter(members, func(m *discordgo.Member, _ int) bool {
		return lo.Contains(names, m.User.Username)
	})

	return result, nil
}
