package webhook

import (
	"context"
	"sort"

	"github.com/bwmarrin/discordgo"
)

// ChannelWebhookIDs returns a list of webhook IDs in a channel.
func ChannelWebhookIDs(ctx context.Context, client *discordgo.Session, channelID string) ([]string, error) {
	webhooks, err := FetchChannelWebhooks(ctx, client, channelID)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, webhook := range webhooks {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			ids = append(ids, webhook.ID)
		}
	}

	// Sort IDs for consistency
	sort.Strings(ids)

	return ids, nil
}

// ChannelWebhookNames returns a list of webhook names in a channel.
func ChannelWebhookNames(ctx context.Context, client *discordgo.Session, channelID string) ([]string, error) {
	webhooks, err := FetchChannelWebhooks(ctx, client, channelID)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, webhook := range webhooks {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			names = append(names, webhook.Name)
		}
	}

	// Sort names for consistency
	sort.Strings(names)

	return names, nil
}

// GuildWebhookIDs returns a list of webhook IDs in a guild.
func GuildWebhookIDs(ctx context.Context, client *discordgo.Session, guildID string) ([]string, error) {
	webhooks, err := FetchGuildWebhooks(ctx, client, guildID)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, webhook := range webhooks {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			ids = append(ids, webhook.ID)
		}
	}

	// Sort IDs for consistency
	sort.Strings(ids)

	return ids, nil
}

// GuildWebhookNames returns a list of webhook names in a guild.
func GuildWebhookNames(ctx context.Context, client *discordgo.Session, guildID string) ([]string, error) {
	webhooks, err := FetchGuildWebhooks(ctx, client, guildID)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, webhook := range webhooks {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			names = append(names, webhook.Name)
		}
	}

	// Sort names for consistency
	sort.Strings(names)

	return names, nil
}
