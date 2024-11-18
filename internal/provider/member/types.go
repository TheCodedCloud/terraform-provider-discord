package member

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MemberDataSource defines the data source implementation.
type MemberDataSource struct {
	client *discordgo.Session
}

// MemberDataSourceModel maps the data source schema data.
type MemberDataSourceModel struct {
	// GuildID is the ID of the guild.
	GuildID types.String `tfsdk:"guild_id"`

	// The ID of the member.
	ID types.String `tfsdk:"id"`

	// The username of the member, not unique across the platform.
	Username types.String `tfsdk:"username"`

	// The user's Discord-tag.
	Discriminator types.String `tfsdk:"discriminator"`

	// The user
	User *User `tfsdk:"user"`

	// This user's guild nickname (if one is set).
	Nick types.String `tfsdk:"nick"`

	// The member's guild avatar hash.
	Avatar types.String `tfsdk:"avatar"`

	// Array of roles
	Roles types.List `tfsdk:"roles"`

	// When the user joined the guild.
	JoinedAt types.String `tfsdk:"joined_at"`

	// When the user started boosting the guild.
	PremiumSince types.String `tfsdk:"premium_since"`

	// Whether the user is deafened in voice channels.
	Deaf types.Bool `tfsdk:"deaf"`

	// Whether the user is muted in voice channels.
	Mute types.Bool `tfsdk:"mute"`

	// Guild member flags
	Flags types.List `tfsdk:"flags"`

	// Whether the user has not yet passed the guild's Membership Screening requirements.
	Pending types.Bool `tfsdk:"pending"`

	// Total permissions of the member in the channel, including overwrites, returned when in the interaction object.
	Permissions types.List `tfsdk:"permissions"`

	// When the user's timeout will expire.
	CommunicationDisabledUntil types.String `tfsdk:"communication_disabled_until"`

	// Data for the member's guild avatar decoration.
	// AvatarDecorationData types.String `tfsdk:"avatar_decoration_data"`
}

// MemberResource defines the resource implementation.
type MemberResource struct {
	client *discordgo.Session
}

// MemberResourceModel maps the resource schema data.
type MemberResourceModel struct {
	// LastUpdated is the last time the resource was updated.
	LastUpdated types.String `tfsdk:"last_updated"`

	MemberDataSourceModel
}
