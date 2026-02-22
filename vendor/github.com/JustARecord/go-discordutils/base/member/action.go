package member

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// AddRole adds a role to a member.
func AddRole(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, member *discordgo.Member, role *discordgo.Role) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return client.GuildMemberRoleAdd(guild.ID, member.User.ID, role.ID)
	}
}

// RemoveRole removes a role from a member.
func RemoveRole(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, member *discordgo.Member, role *discordgo.Role) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return client.GuildMemberRoleRemove(guild.ID, member.User.ID, role.ID)
	}
}
