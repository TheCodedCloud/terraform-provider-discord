package role

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleDataSource defines the data source implementation.
type RoleDataSource struct {
	client *discordgo.Session
}

// RoleDataSourceModel maps the data source schema data.
type RoleDataSourceModel struct {
	// GuildID is the ID of the guild.
	GuildID types.String `tfsdk:"guild_id"`

	// The ID of the role.
	ID types.String `tfsdk:"id"`

	// The name of the role.
	Name types.String `tfsdk:"name"`

	// The hex color of this role.
	Color types.String `tfsdk:"color"`

	// Whether this role is hoisted (shows up separately in member list).
	Hoist types.Bool `tfsdk:"hoist"`

	// The hash of the role icon. Use Role.IconURL to retrieve the icon's URL.
	Icon types.String `tfsdk:"icon"`

	// The emoji assigned to this role.
	UnicodeEmoji types.String `tfsdk:"unicode_emoji"`

	// The position of this role in the guild's role hierarchy.
	Position types.Int32 `tfsdk:"position"`

	// The permissions of the role on the guild (doesn't include channel overrides).
	// This is a combination of bit masks; the presence of a certain permission can
	// be checked by performing a bitwise AND between this int and the permission.
	Permissions types.List `tfsdk:"permissions"`

	// Whether this role is managed by an integration, and
	// thus cannot be manually added to, or taken from, members.
	Managed types.Bool `tfsdk:"managed"`

	// Whether this role is mentionable.
	Mentionable types.Bool `tfsdk:"mentionable"`

	// The tags this role has
	// Tags RoleTags `tfsdk:"tags"`

	// The flags of the role, which describe its extra features.
	// This is a combination of bit masks; the presence of a certain flag can
	// be checked by performing a bitwise AND between this int and the flag.
	// Flags RoleFlags `json:"flags"`
	Flags types.List `tfsdk:"flags"`
}

// RoleResource defines the resource implementation.
type RoleResource struct {
	client *discordgo.Session
}

// RoleResourceModel maps the resource schema data.
type RoleResourceModel struct {
	// LastUpdated is the last time the resource was updated.
	LastUpdated types.String `tfsdk:"last_updated"`

	RoleDataSourceModel
}
