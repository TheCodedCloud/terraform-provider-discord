package role

import (
	"context"
	"fmt"

	"github.com/JustARecord/go-discordutils/base/member"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

// FetchByName fetches a role by name.
func FetchByName(ctx context.Context, client *discordgo.Session, guild, name string) (*discordgo.Role, error) {
	roles, err := AllByID(ctx, client, guild)
	if err != nil {
		return nil, err
	}

	role, ok := lo.Find(roles, func(r *discordgo.Role) bool {
		return r.Name == name
	})

	if !ok {
		return nil, fmt.Errorf("role not found: name=%s", name)
	}

	return role, nil
}

// FetchByID fetches a role by ID.
func FetchByID(ctx context.Context, client *discordgo.Session, guild, id string) (*discordgo.Role, error) {
	roles, err := AllByID(ctx, client, guild)
	if err != nil {
		return nil, err
	}

	role, ok := lo.Find(roles, func(r *discordgo.Role) bool {
		return r.ID == id
	})

	if !ok {
		return nil, fmt.Errorf("role not found: id=%s", id)
	}

	return role, nil
}

// FetchByNames fetches roles by their names.
func FetchByNames(ctx context.Context, client *discordgo.Session, guild string, names []string) ([]*discordgo.Role, error) {
	roles, err := AllByID(ctx, client, guild)
	if err != nil {
		return nil, err
	}

	result := lo.Filter(roles, func(r *discordgo.Role, _ int) bool {
		return lo.Contains(names, r.Name)
	})

	return result, nil
}

// FetchByIDs fetches roles by their IDs.
func FetchByIDs(ctx context.Context, client *discordgo.Session, guild string, ids []string) ([]*discordgo.Role, error) {
	roles, err := AllByID(ctx, client, guild)
	if err != nil {
		return nil, err
	}

	result := lo.Filter(roles, func(r *discordgo.Role, _ int) bool {
		return lo.Contains(ids, r.ID)
	})

	return result, nil
}

// FetchMembers fetches members with a role.
func FetchMembers(ctx context.Context, client *discordgo.Session, guild_id, role_id string) ([]*discordgo.Member, error) {
	// Fetch all members in the guild
	members, err := member.Fetch(ctx, client, guild_id, -1)
	if err != nil {
		return nil, err
	}

	result := lo.Filter(members, func(m *discordgo.Member, _ int) bool {
		return lo.Contains(m.Roles, role_id)
	})

	return result, nil
}

// AllByID fetches all roles in a guild by ID.
func AllByID(ctx context.Context, client *discordgo.Session, guild string) ([]*discordgo.Role, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.GuildRoles(guild)
	}
}

// All fetches all roles in a guild.
func All(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild) ([]*discordgo.Role, error) {
	return AllByID(ctx, client, guild.ID)
}

// AllMap fetches all roles in a guild as a map with the role name as the key.
func AllMap(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild) (map[string]*discordgo.Role, error) {
	roles, err := All(ctx, client, guild)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*discordgo.Role)
	for _, r := range roles {
		result[r.Name] = r
	}

	return result, nil
}
