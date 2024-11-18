package role_members

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleMembersDataSource defines the data source implementation.
type RoleMembersDataSource struct {
	client *discordgo.Session
}

// RoleMembersDataSourceModel maps the data source schema data.
type RoleMembersDataSourceModel struct {
	// GuildID is the ID of the guild.
	GuildID types.String `tfsdk:"guild_id"`

	// The role ID.
	RoleID types.String `tfsdk:"role_id"`

	// The name of the role.
	Role types.String `tfsdk:"role"`

	// Array of role members.
	Members types.List `tfsdk:"members"`
}

// RoleMembersResource defines the resource implementation.
type RoleMembersResource struct {
	client *discordgo.Session
}

// RoleMembersResourceModel maps the resource schema data.
type RoleMembersResourceModel struct {
	// LastUpdated is the last time the resource was updated.
	LastUpdated types.String `tfsdk:"last_updated"`

	RoleMembersDataSourceModel
}
