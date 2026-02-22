package role

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/JustARecord/go-discordutils/base/member"
	"github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/schollz/progressbar/v3"
)

// Create creates a role.
func Create(ctx context.Context, client *discordgo.Session, guild_id string, params *discordgo.RoleParams) (*discordgo.Role, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.GuildRoleCreate(guild_id, params)
	}
}

// UpdateByID updates a role by ID.
func UpdateByID(ctx context.Context, client *discordgo.Session, guild_id, role_id string, params *discordgo.RoleParams) (*discordgo.Role, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.GuildRoleEdit(guild_id, role_id, params)
	}
}

// Update updates a role.
func Update(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, role *discordgo.Role, params *discordgo.RoleParams) (*discordgo.Role, error) {
	return UpdateByID(ctx, client, guild.ID, role.ID, params)
}

// UpdateByName updates a role by name.
func UpdateByName(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, name string, params *discordgo.RoleParams) (*discordgo.Role, error) {
	role, err := FetchByName(ctx, client, guild.ID, name)
	if err != nil {
		return nil, err
	}

	return Update(ctx, client, guild, role, params)
}

// DeleteByID deletes a role by ID.
func DeleteByID(ctx context.Context, client *discordgo.Session, guild_id, role_id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return client.GuildRoleDelete(guild_id, role_id)
	}
}

// Delete deletes a role.
func Delete(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, role *discordgo.Role) error {
	return DeleteByID(ctx, client, guild.ID, role.ID)
}

// DeleteByName deletes a role by name.
func DeleteByName(ctx context.Context, client *discordgo.Session, guild_id, name string) error {
	role, err := FetchByName(ctx, client, guild_id, name)
	if err != nil {
		return err
	}

	return DeleteByID(ctx, client, guild_id, role.ID)
}

// AddMembers adds members to a role.
func AddMembers(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, role *discordgo.Role, members []*discordgo.Member) error {
	bar := progressbar.NewOptions(
		len(members),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetDescription(fmt.Sprintf("adding %s role", role.Name)),
	)

	// Add the role to all members
	for _, m := range members {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := member.AddRole(ctx, client, guild, m, role); err != nil {
				return err
			}

			bar.Add(1)
			slog.Debug("Added role to member", "role", role.Name, "member", m.User.Username)

			// Sleep for a small duration to avoid hitting rate limits
			duration := time.Duration(100) * time.Millisecond
			utils.Sleep(ctx, duration)
		}
	}

	bar.Finish()

	return nil
}

// SetMembers sets the members of a role.
func SetMembers(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, role *discordgo.Role, new []*discordgo.Member) ([]*discordgo.Member, error) {
	// TODO: Optimize this function

	// 1. Fetch the current members of the role
	current, err := FetchMembers(ctx, client, guild.ID, role.ID)
	if err != nil {
		return nil, err
	}

	// 2. Remove the members that are not in the new list
	toRemove, remaining := lo.FilterReject(current, func(m *discordgo.Member, _ int) bool {
		_, ok := lo.Find(new, func(a *discordgo.Member) bool {
			return a.User.ID == m.User.ID
		})
		return !ok
	})

	// 3. Add the members that are not in the current list
	toAdd, _ := lo.FilterReject(new, func(m *discordgo.Member, _ int) bool {
		_, ok := lo.Find(remaining, func(a *discordgo.Member) bool {
			return a.User.ID == m.User.ID
		})
		return !ok
	})

	// 4. Remove the role from the members that have it
	if err := RemoveMembers(ctx, client, guild, role, toRemove); err != nil {
		return nil, err
	}

	// 5. Add the role to the members that don't have it
	if err := AddMembers(ctx, client, guild, role, toAdd); err != nil {
		return nil, err
	}

	// 6. Refresh the members that have the role
	result, err := FetchMembers(ctx, client, guild.ID, role.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RemoveMembers removes members from a role.
func RemoveMembers(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, role *discordgo.Role, members []*discordgo.Member) error {
	bar := progressbar.NewOptions(
		len(members),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetDescription(fmt.Sprintf("removing %s role", role.Name)),
	)

	// Remove the role from all members
	for _, m := range members {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := member.RemoveRole(ctx, client, guild, m, role); err != nil {
				return err
			}

			bar.Add(1)
			slog.Debug("Removed role from member", "role", role.Name, "member", m.User.Username)

			// Sleep for a small duration to avoid hitting rate limits
			duration := time.Duration(100) * time.Millisecond
			utils.Sleep(ctx, duration)
		}
	}

	bar.Finish()

	return nil
}
