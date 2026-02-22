package webhook

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

// FetchByID fetches a webhook by ID.
func FetchByID(ctx context.Context, client *discordgo.Session, id string) (*discordgo.Webhook, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.Webhook(id)
	}
}

// ChannelWebhooks fetches all webhooks in a channel.
func FetchChannelWebhooks(ctx context.Context, client *discordgo.Session, channelID string) ([]*discordgo.Webhook, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.ChannelWebhooks(channelID)
	}
}

// FetchChannelWebhookByName fetches a channel webhook by its name.
func FetchChannelWebhookByName(ctx context.Context, client *discordgo.Session, channelID, name string) (*discordgo.Webhook, error) {
	webhooks, err := FetchChannelWebhooks(ctx, client, channelID)
	if err != nil {
		return nil, err
	}

	webhook, ok := lo.Find(webhooks, func(w *discordgo.Webhook) bool {
		return w.Name == name
	})

	if !ok {
		return nil, nil
	}

	return webhook, nil
}

// FetchChannelWebhookByID fetches a channel webhook by its ID.
func FetchChannelWebhookByID(ctx context.Context, client *discordgo.Session, channelID, webhookID string) (*discordgo.Webhook, error) {
	webhook, err := FetchByID(ctx, client, webhookID)
	if err != nil {
		return nil, err
	}

	if webhook.ChannelID != channelID {
		return nil, fmt.Errorf("webhook not found in channel: id=%s", channelID)
	}

	return webhook, nil
}

// GuildWebhooks fetches all webhooks in a guild.
func FetchGuildWebhooks(ctx context.Context, client *discordgo.Session, guildID string) ([]*discordgo.Webhook, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.GuildWebhooks(guildID)
	}
}

// FetchGuildWebhookByName fetches a guild webhook by its name.
func FetchGuildWebhookByName(ctx context.Context, client *discordgo.Session, guildID, name string) (*discordgo.Webhook, error) {
	webhooks, err := FetchGuildWebhooks(ctx, client, guildID)
	if err != nil {
		return nil, err
	}

	webhook, ok := lo.Find(webhooks, func(w *discordgo.Webhook) bool {
		return w.Name == name
	})

	if !ok {
		return nil, nil
	}

	return webhook, nil
}

// FetchGuildWebhookByID fetches a guild webhook by its ID.
func FetchGuildWebhookByID(ctx context.Context, client *discordgo.Session, guildID, webhookID string) (*discordgo.Webhook, error) {
	webhook, err := FetchByID(ctx, client, webhookID)
	if err != nil {
		return nil, err
	}

	if webhook.GuildID != guildID {
		return nil, fmt.Errorf("webhook not found in guild: id=%s", guildID)
	}

	return webhook, nil
}

// FetchIntegrationWebhooks fetches all webhooks for an integration.
func FetchIntegrationWebhooks(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, integration *discordgo.Integration) ([]*discordgo.Webhook, error) {
	guildWebhooks, err := FetchGuildWebhooks(ctx, client, guild.ID)
	if err != nil {
		return nil, err
	}

	result := lo.Filter(guildWebhooks, func(w *discordgo.Webhook, _ int) bool {
		return w.ApplicationID == integration.Account.ID
	})

	return result, nil
}
