package discord

import (
	"context"
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
)

// CalcPermissions calculates the permissions from a list of permissions.
func CalcPermissions(permissions []string) int64 {
	var perms int64
	for _, perm := range permissions {
		perms |= Permissions[perm]
	}
	return perms
}

// ParseOverwrite parses the provided permission overwrite into allowed and denied permissions lists.
func ParseOverwrite(overwrite *discordgo.PermissionOverwrite) ([]string, []string, error) {
	allowed := ListStringify(overwrite.Allow)
	denied := ListStringify(overwrite.Deny)

	sort.Strings(allowed)
	sort.Strings(denied)

	return allowed, denied, nil
}

// FetchChannelPermissions fetches the permissions for the provided channel.
func FetchChannelPermissions(ctx context.Context, client *discordgo.Session, guild_id, channel_id, id, permissionType string) (*discordgo.PermissionOverwrite, error) {
	// Fetch the channel
	channel, err := FetchChannelByID(ctx, client, guild_id, channel_id)
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
	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.Type == permissionOverwriteType && overwrite.ID == id {
			return overwrite, nil
		}
	}

	return nil, fmt.Errorf("permission overwrite not found for %s type %s", permissionType, id)
}

// CreatePermissionOverwrite creates a new permission overwrite for the provided channel.
func CreatePermissionOverwrite(ctx context.Context, client *discordgo.Session, guild_id, channel_id, id, permissionType string, allow, deny []string) (*discordgo.PermissionOverwrite, error) {
	permissionOverwriteType := discordgo.PermissionOverwriteTypeRole
	if permissionType == "member" {
		permissionOverwriteType = discordgo.PermissionOverwriteTypeMember
	}

	allowInt := CalcPermissions(allow)
	denyInt := CalcPermissions(deny)

	// Create the permission overwrite
	err := client.ChannelPermissionSet(channel_id, id, permissionOverwriteType, allowInt, denyInt)
	if err != nil {
		return nil, err
	}

	// Fetch the overwrite to ensure it was created
	return FetchChannelPermissions(ctx, client, guild_id, channel_id, id, permissionType)
}

// UpdatePermissionOverwrite updates the provided permission overwrite with the new permissions.
// This is just an alias for CreatePermissionOverwrite as that function internally calls ChannelPermissionSet which creates or updates the permission overwrite.
func UpdatePermissionOverwrite(ctx context.Context, client *discordgo.Session, guild_id, channel_id, id, permissionType string, allow, deny []string) (*discordgo.PermissionOverwrite, error) {
	return CreatePermissionOverwrite(ctx, client, guild_id, channel_id, id, permissionType, allow, deny)
}

// DeletePermissionOverwrite deletes the provided permission overwrite.
func DeletePermissionOverwrite(ctx context.Context, client *discordgo.Session, guild_id, channel_id, id, permissionType string) error {
	return client.ChannelPermissionDelete(channel_id, id)
}
