package guild

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GuildDataSource defines the data source implementation.
type GuildDataSource struct {
	client *discordgo.Session
}

// GuildDataSourceModel maps the data source schema data.
type GuildDataSourceModel struct {
	// The ID of the guild.
	ID types.String `tfsdk:"id"`

	// The name of the guild. (2–100 characters)
	Name types.String `tfsdk:"name"`

	// The hash of the guild's icon. Use Session.GuildIcon
	// to retrieve the icon itself.
	Icon types.String `tfsdk:"icon"`

	// The hash of the guild's splash.
	Splash types.String `tfsdk:"splash"`

	// The hash of the guild's discovery splash.
	DiscoverySplash types.String `tfsdk:"discovery_splash"`

	// If we are the owner of the guild
	Owner types.Bool `tfsdk:"owner"`

	// The user ID of the owner of the guild.
	OwnerID types.String `tfsdk:"owner_id"`

	// The voice region of the guild.
	// deprecated:
	Region types.String `tfsdk:"region"`

	// The ID of the AFK voice channel.
	AfkChannelID types.String `tfsdk:"afk_channel_id"`

	// The timeout, in seconds, before a user is considered AFK in voice.
	AfkTimeout types.Int32 `tfsdk:"afk_timeout"`

	// Whether or not the Server Widget is enabled
	WidgetEnabled types.Bool `tfsdk:"widget_enabled"`

	// The Channel ID for the Server Widget
	WidgetChannelID types.String `tfsdk:"widget_channel_id"`

	// The verification level required for the guild.
	VerificationLevel types.String `tfsdk:"verification_level"`

	// The default message notification setting for the guild.
	DefaultMessageNotifications types.String `tfsdk:"default_message_notifications"`

	// The explicit content filter level
	ExplicitContentFilter types.String `tfsdk:"explicit_content_filter"`

	// Enabled guild features
	Features types.List `tfsdk:"features"`

	// Required MFA level for the guild
	MfaLevel types.String `tfsdk:"mfa_level"`

	// The application id of the guild if bot created.
	ApplicationID types.String `tfsdk:"application_id"`

	// The Channel ID to which system messages are sent (eg join and leave messages)
	SystemChannelID types.String `tfsdk:"system_channel_id"`

	// The System channel flags
	SystemChannelFlags types.List `tfsdk:"system_channel_flags"`

	// The ID of the rules channel ID, used for rules.
	RulesChannelID types.String `tfsdk:"rules_channel_id"`

	// The maximum number of presences for the guild (the default value, currently 25000, is in effect when null is returned)
	MaxPresences types.Int32 `tfsdk:"max_presences"`

	// The maximum number of members for the guild
	MaxMembers types.Int32 `tfsdk:"max_members"`

	// the vanity url code for the guild
	VanityURLCode types.String `tfsdk:"vanity_url_code"`

	// the description for the guild
	Description types.String `tfsdk:"description"`

	// The hash of the guild's banner
	Banner types.String `tfsdk:"banner"`

	// The premium tier of the guild (Server Boost Level)
	PremiumTier types.String `tfsdk:"premium_tier"`

	// The total number of users currently boosting this server
	PremiumSubscriptionCount types.Int32 `tfsdk:"premium_subscription_count"`

	// The preferred locale of a guild with the "PUBLIC" feature; used in server discovery and notices from Discord; defaults to "en-US"
	PreferredLocale types.String `tfsdk:"preferred_locale"`

	// The id of the channel where admins and moderators of guilds with the "PUBLIC" feature receive notices from Discord
	PublicUpdatesChannelID types.String `tfsdk:"public_updates_channel_id"`

	// The maximum amount of users in a video channel
	MaxVideoChannelUsers types.Int32 `tfsdk:"max_video_channel_users"`

	// The maximum amount of users in a stage video channel
	MaxStageVideoChannelUsers types.Int32 `tfsdk:"max_stage_video_channel_users"`

	// Approximate number of members in this guild, returned from the GET /guild/<id> endpoint when with_counts is true
	ApproximateMemberCount types.Int32 `tfsdk:"approximate_member_count"`

	// Approximate number of non-offline members in this guild, returned from the GET /guild/<id> endpoint when with_counts is true
	ApproximatePresenceCount types.Int32 `tfsdk:"approximate_presence_count"`

	// The NSFW Level of the guild
	NSFWLevel types.String `tfsdk:"nsfw_level"`

	// whether the guild has the boost progress bar enabled
	PremiumProgressBarEnabled types.Bool `tfsdk:"premium_progress_bar_enabled"`

	// The id of the channel where admins and moderators of Community guilds receive safety alerts from Discord
	SafetyAlertsChannelID types.String `tfsdk:"safety_alerts_channel_id"`

	// Permissions of our user
	Permissions types.Int64 `tfsdk:"permissions"`
}
