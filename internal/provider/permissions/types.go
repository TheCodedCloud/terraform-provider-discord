package permissions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PermissionsDataSource defines the data source implementation.
type PermissionsDataSource struct {
	client *discordgo.Session
}

// PermissionsDataSourceModel maps the data source schema data.
type PermissionsDataSourceModel struct {
	// GuildID is the ID of the guild.
	GuildID types.String `tfsdk:"guild_id"`

	// ChannelID is the ID of the channel.
	ChannelID types.String `tfsdk:"channel_id"`

	// The ID of the role or user.
	ID types.String `tfsdk:"id"`

	// The type of the overwrite, either "role" or "member".
	Type types.String `tfsdk:"type"`

	// The permissions that are allowed.
	Allow types.List `tfsdk:"allow"`

	// The permissions that are denied.
	Deny types.List `tfsdk:"deny"`
}

// PermissionsResource defines the resource implementation.
type PermissionsResource struct {
	client *discordgo.Session
}

// PermissionsResourceModel maps the resource schema data.
type PermissionsResourceModel struct {
	// LastUpdated is the last time the resource was updated.
	LastUpdated types.String `tfsdk:"last_updated"`

	PermissionsDataSourceModel
}
