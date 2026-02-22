package webhook

import (
	"context"
	"log/slog"

	"github.com/JustARecord/go-discordutils/base/channel"
	"github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

// FetchMessageID fetches a webhook message by ID.
func FetchMessageID(ctx context.Context, client *discordgo.Session, webhookID, webhookToken, messageID string) (*discordgo.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.WebhookMessage(webhookID, webhookToken, messageID)
	}
}

// FetchMessage fetches a webhook message.
func FetchMessage(ctx context.Context, client *discordgo.Session, webhook *discordgo.Webhook, messageID string) (*discordgo.Message, error) {
	return FetchMessageID(ctx, client, webhook.ID, webhook.Token, messageID)
}

// RefreshMessages refreshes messages with webhook data.
func RefreshMessages(ctx context.Context, client *discordgo.Session, webhook *discordgo.Webhook, messages []*discordgo.Message) ([]*discordgo.Message, error) {
	if webhook == nil {
		return messages, nil
	}

	result := []*discordgo.Message{}

	matchingMessages := lo.Filter(messages, func(message *discordgo.Message, _ int) bool {
		return message.WebhookID != "" && message.WebhookID == webhook.ID
	})

	for _, message := range matchingMessages {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			slog.Info("Refreshing message", "message", message, "webhook", message.WebhookID, "expected", webhook.ID)

			refreshedMessage, err := FetchMessage(ctx, client, webhook, message.ID)
			if err != nil {
				// If the message is not found, skip it
				if utils.NotFoundError(err) {
					slog.Warn("Message not found", "message", message.ID, "webhook", message.WebhookID, "expected", webhook.ID)
					continue
				}

				slog.Error("Error fetching message", "message", message.ID, "error", err)
				return nil, err
			}

			result = append(result, refreshedMessage)
		}
	}

	return result, nil
}

// FetchChannelMessages fetches messages in a channel.
func FetchChannelMessages(ctx context.Context, client *discordgo.Session, ch *discordgo.Channel, wh *discordgo.Webhook, limit int) ([]*discordgo.Message, error) {
	messages, err := channel.FetchMessages(ctx, client, ch, limit)
	if err != nil {
		return nil, err
	}

	slog.Info("Fetched messages", "channel", ch.Name, "count", len(messages))

	messages, err = RefreshMessages(ctx, client, wh, messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// EditMessageID edits a webhook message.
func EditMessageID(ctx context.Context, client *discordgo.Session, webhook *discordgo.Webhook, messageID string, params *discordgo.WebhookEdit, thread *discordgo.Channel) (*discordgo.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if thread != nil {
			return EditThreadMessage(ctx, client, webhook, thread, messageID, params)
		} else {
			return client.WebhookMessageEdit(webhook.ID, webhook.Token, messageID, params)
		}
	}
}

// EditMessage edits a webhook message.
func EditMessage(ctx context.Context, client *discordgo.Session, webhook *discordgo.Webhook, message *discordgo.Message, params *discordgo.WebhookEdit, thread *discordgo.Channel) (*discordgo.Message, error) {
	return EditMessageID(ctx, client, webhook, message.ID, params, thread)
}

// DeleteMessageID deletes a webhook message.
func DeleteMessageID(ctx context.Context, client *discordgo.Session, webhookID, webhookToken, messageID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return client.WebhookMessageDelete(webhookID, webhookToken, messageID)
	}
}

// DeleteMessage deletes a webhook message.
func DeleteMessage(ctx context.Context, client *discordgo.Session, webhook *discordgo.Webhook, message *discordgo.Message) error {
	// TODO: support threads
	return DeleteMessageID(ctx, client, webhook.ID, webhook.Token, message.ID)
}

// DeleteMessages deletes multiple webhook messages.
func DeleteMessages(ctx context.Context, client *discordgo.Session, webhook *discordgo.Webhook, messages []*discordgo.Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		for _, message := range messages {
			if err := DeleteMessage(ctx, client, webhook, message); err != nil {
				return err
			}
		}

		return nil
	}
}
