package webhook

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WebhookDataSource defines the data source implementation.
type WebhookDataSource struct {
	client *discordgo.Session
}

// WebhookDataSourceModel maps the data source schema data.
type WebhookDataSourceModel struct {
	// The ID of the webhook.
	ID types.String `tfsdk:"id"`

	// The type of the webhook, either "INCOMING", "CHANNEL_FOLLOWER", or "APPLICATION".
	Type types.String `tfsdk:"type"`

	// The guild id this webhook is for, if any.
	GuildID types.String `tfsdk:"guild_id"`

	// The channel id this webhook is for, if any.
	ChannelID types.String `tfsdk:"channel_id"`

	// The default name of the webhook.
	Name types.String `tfsdk:"name"`

	// The default user avatar hash of the webhook.
	Avatar types.String `tfsdk:"avatar"`

	// The secure token of the webhook (returned for Incoming Webhooks).
	Token types.String `tfsdk:"token"`

	// The Bot/OAuth2 application that created this webhook.
	ApplicationID types.String `tfsdk:"application_id"`

	// The partial guild of the channel that this webhook is following (returned for Channel Follower Webhooks).
	// SourceGuild types.String `tfsdk:"source_guild"`

	// The partial channel that this webhook is following (returned for Channel Follower Webhooks).
	// SourceChannel types.String `tfsdk:"source_channel"`
}

// WebhookResource defines the resource implementation.
type WebhookResource struct {
	client *discordgo.Session
}

// WebhookResourceModel maps the resource schema data.
type WebhookResourceModel struct {
	// LastUpdated is the last time the resource was updated.
	LastUpdated types.String `tfsdk:"last_updated"`

	WebhookDataSourceModel
}
