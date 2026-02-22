package webhook

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// CreateWebhook creates a webhook in a channel.
func CreateWebhook(ctx context.Context, client *discordgo.Session, channelID, name, avatar string) (*discordgo.Webhook, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.WebhookCreate(channelID, name, avatar)
	}
}

// UpdateWebhook updates a webhook.
func UpdateWebhook(ctx context.Context, client *discordgo.Session, webhookID, name, avatar, channelID string) (*discordgo.Webhook, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.WebhookEdit(webhookID, name, avatar, channelID)
	}
}

// DeleteWebhook deletes a webhook.
func DeleteWebhook(ctx context.Context, client *discordgo.Session, webhookID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return client.WebhookDelete(webhookID)
	}
}
