package permissions

import (
	"context"

	"github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
)

// CreatePermissionOverwrite creates a new permission overwrite for the provided channel.
func CreatePermissionOverwrite(ctx context.Context, client *discordgo.Session, guild_id, channel_id, id, permissionType string, allow, deny []string) (*discordgo.PermissionOverwrite, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		permissionOverwriteType := discordgo.PermissionOverwriteTypeRole
		if permissionType == "member" {
			permissionOverwriteType = discordgo.PermissionOverwriteTypeMember
		}

		allowInt := utils.CalcPermissions(allow)
		denyInt := utils.CalcPermissions(deny)

		// Create the permission overwrite
		err := client.ChannelPermissionSet(channel_id, id, permissionOverwriteType, allowInt, denyInt)
		if err != nil {
			return nil, err
		}

		// Fetch the overwrite to ensure it was created
		return FetchChannelPermissions(ctx, client, guild_id, channel_id, id, permissionType)
	}
}

// UpdatePermissionOverwrite updates the provided permission overwrite with the new permissions.
// This is just an alias for CreatePermissionOverwrite as that function internally calls ChannelPermissionSet which creates or updates the permission overwrite.
func UpdatePermissionOverwrite(ctx context.Context, client *discordgo.Session, guild_id, channel_id, id, permissionType string, allow, deny []string) (*discordgo.PermissionOverwrite, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return CreatePermissionOverwrite(ctx, client, guild_id, channel_id, id, permissionType, allow, deny)
	}
}

// DeletePermissionOverwrite deletes the provided permission overwrite.
func DeletePermissionOverwrite(ctx context.Context, client *discordgo.Session, guild_id, channel_id, id, permissionType string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return client.ChannelPermissionDelete(channel_id, id)
	}
}
