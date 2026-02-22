package member

import (
	"context"
	"slices"
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

// Names returns a list of member names.
func Names(members []*discordgo.Member) []string {
	names := []string{}
	for _, m := range members {
		names = append(names, m.User.Username)
	}

	// Sort the names
	sort.Strings(names)

	return names
}

// IDs returns a list of member IDs.
func IDs(members []*discordgo.Member) []string {
	ids := []string{}
	for _, m := range members {
		ids = append(ids, m.User.ID)
	}

	// Sort the IDs
	sort.Strings(ids)

	return ids
}

// HasRole checks if a member has a role.
func HasRole(ctx context.Context, member *discordgo.Member, role *discordgo.Role) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		return slices.Contains(member.Roles, role.ID), nil
	}
}

// HasRoles checks if a member has all the roles.
func HasRoles(ctx context.Context, member *discordgo.Member, roles []*discordgo.Role) ([]*discordgo.Role, error) {
	var err error

	result := lo.Filter(roles, func(r *discordgo.Role, _ int) bool {
		has, hasErr := HasRole(ctx, member, r)
		if hasErr != nil {
			err = hasErr
			return false
		}

		return has
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
