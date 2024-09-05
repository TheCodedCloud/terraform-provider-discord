package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// FetchWebhookByID fetches a webhook by its ID.
func FetchWebhookByID(ctx context.Context, client *discordgo.Session, id string) (*discordgo.Webhook, error) {
	return client.Webhook(id)
}

// ChannelWebhooks fetches all webhooks in a channel.
func ChannelWebhooks(ctx context.Context, client *discordgo.Session, channelID string) ([]*discordgo.Webhook, error) {
	return client.ChannelWebhooks(channelID)
}

// FetchChannelWebhookByName fetches a channel webhook by its name.
func FetchChannelWebhookByName(ctx context.Context, client *discordgo.Session, channelID, name string) (*discordgo.Webhook, error) {
	webhooks, err := ChannelWebhooks(ctx, client, channelID)
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Name == name {
			return webhook, nil
		}
	}

	return nil, nil
}

// GuildWebhooks fetches all webhooks in a guild.
func GuildWebhooks(ctx context.Context, client *discordgo.Session, guildID string) ([]*discordgo.Webhook, error) {
	return client.GuildWebhooks(guildID)
}

// FetchGuildWebhookByName fetches a guild webhook by its name.
func FetchGuildWebhookByName(ctx context.Context, client *discordgo.Session, guildID, name string) (*discordgo.Webhook, error) {
	webhooks, err := GuildWebhooks(ctx, client, guildID)
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Name == name {
			return webhook, nil
		}
	}

	return nil, nil
}

// CreateWebhook creates a webhook in a channel.
func CreateWebhook(ctx context.Context, client *discordgo.Session, channelID, name, avatar string) (*discordgo.Webhook, error) {
	return client.WebhookCreate(channelID, name, avatar)
}

// UpdateWebhook updates a webhook.
func UpdateWebhook(ctx context.Context, client *discordgo.Session, webhookID, name, avatar, channelID string) (*discordgo.Webhook, error) {
	return client.WebhookEdit(webhookID, name, avatar, channelID)
}

// DeleteWebhook deletes a webhook.
func DeleteWebhook(ctx context.Context, client *discordgo.Session, webhookID string) error {
	return client.WebhookDelete(webhookID)
}
