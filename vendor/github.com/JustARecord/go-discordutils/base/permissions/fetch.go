package permissions

import (
	"context"
	"fmt"

	"github.com/JustARecord/go-discordutils/base/channel"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

// FetchChannelPermissions fetches the permissions for the provided channel.
func FetchChannelPermissions(ctx context.Context, client *discordgo.Session, guild_id, channel_id, id, permissionType string) (*discordgo.PermissionOverwrite, error) {
	// Fetch the channel
	channel, err := channel.FetchByID(ctx, client, guild_id, channel_id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channel from guild %s with id %s: %w", guild_id, channel_id, err)
	}

	var permissionOverwriteType discordgo.PermissionOverwriteType
	if permissionType == "role" {
		permissionOverwriteType = discordgo.PermissionOverwriteTypeRole
	} else if permissionType == "member" {
		permissionOverwriteType = discordgo.PermissionOverwriteTypeMember
	} else {
		return nil, fmt.Errorf("invalid permission type: %s", permissionType)
	}

	// Find the permission overwrite
	overwrite, ok := lo.Find(channel.PermissionOverwrites, func(overwrite *discordgo.PermissionOverwrite) bool {
		return overwrite.Type == permissionOverwriteType && overwrite.ID == id
	})

	if !ok {
		return nil, fmt.Errorf("permission overwrite not found for %s type %s", permissionType, id)
	}

	return overwrite, nil
}
