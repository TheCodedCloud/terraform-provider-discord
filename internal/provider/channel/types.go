package channel

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ChannelDataSource defines the data source implementation.
type ChannelDataSource struct {
	client *discordgo.Session
}

// ChannelDataSourceModel maps the data source schema data.
type ChannelDataSourceModel struct {
	// The ID of the channel.
	ID types.String `tfsdk:"id"`

	// The type of the channel.
	Type types.String `tfsdk:"type"`

	// The ID of the guild to which the channel belongs, if it is in a guild.
	// Else, this ID is empty (e.g. DM channels).
	GuildID types.String `tfsdk:"guild_id"`

	// The position of the channel, used for sorting in client.
	Position types.Int32 `tfsdk:"position"`

	// The name of the channel.
	Name types.String `tfsdk:"name"`

	// The topic of the channel.
	Topic types.String `tfsdk:"topic"`

	// Whether the channel is marked as NSFW.
	NSFW types.Bool `tfsdk:"nsfw"`

	// The bitrate of the channel, if it is a voice channel.
	Bitrate types.Int32 `tfsdk:"bitrate"`

	// The user limit of the voice channel.
	UserLimit types.Int32 `tfsdk:"user_limit"`

	// Amount of seconds a user has to wait before sending another message or creating another thread (0-21600)
	// bots, as well as users with the permission manage_messages or manage_channel, are unaffected
	RateLimitPerUser types.Int32 `tfsdk:"rate_limit_per_user"`

	// The recipients of the channel. This is only populated in DM channels.
	// Recipients []*User `tfsdk:"recipients"`

	// Icon of the group DM channel.
	Icon types.String `tfsdk:"icon"`

	// ID of the creator of the group DM or thread
	OwnerID types.String `tfsdk:"owner_id"`

	// ApplicationID of the DM creator Zeroed if guild channel or not a bot user
	ApplicationID types.String `tfsdk:"application_id"`

	// for group DM channels: whether the channel is managed by an application
	// via the gdm.join OAuth2 scope
	// Managed types.Bool `tfsdk:"managed"`

	// The ID of the parent channel, if the channel is under a category. For threads - id of the channel thread was created in.
	ParentID types.String `tfsdk:"parent_id"`

	// The IDs of the child channels of the category, if the channel is a category.
	Children types.List `tfsdk:"children"`

	// The timestamp of the last pinned message in the channel.
	// nil if the channel has no pinned messages.
	LastPinTimestamp types.String `tfsdk:"last_pin_timestamp"`

	// Thread-specific fields not needed by other channels
	// ThreadMetadata types.Object `tfsdk:"thread_metadata"`

	// Thread member object for the current user, if they have joined the thread, only included on certain API endpoints
	// ThreadMember *ThreadMember `tfsdk:"thread_member"`

	// Channel flags.
	Flags types.List `tfsdk:"flags"`

	// The set of tags that can be used in a forum channel.
	// AvailableTags types.List `tfsdk:"available_tags"`

	// The IDs of the set of tags that have been applied to a thread in a forum channel.
	AppliedTags types.List `tfsdk:"applied_tags"`

	// Emoji to use as the default reaction to a forum post.
	// DefaultReactionEmoji ForumDefaultReaction `tfsdk:"default_reaction_emoji"`

	// The initial RateLimitPerUser to set on newly created threads in a channel.
	// This field is copied to the thread at creation time and does not live update.
	DefaultThreadRateLimitPerUser types.Int32 `tfsdk:"default_thread_rate_limit_per_user"`

	// The default sort order type used to order posts in forum channels.
	// Defaults to null, which indicates a preferred sort order hasn't been set by a channel admin.
	DefaultSortOrder types.String `tfsdk:"default_sort_order"`

	// The default forum layout view used to display posts in forum channels.
	// Defaults to ForumLayoutNotSet, which indicates a layout view has not been set by a channel admin.
	DefaultForumLayout types.String `tfsdk:"default_forum_layout"`
}

// ChannelResource defines the resource implementation.
type ChannelResource struct {
	client *discordgo.Session
}

// ChannelResourceModel maps the resource schema data.
type ChannelResourceModel struct {
	// LastUpdated is the last time the resource was updated.
	LastUpdated types.String `tfsdk:"last_updated"`

	ChannelDataSourceModel
}
