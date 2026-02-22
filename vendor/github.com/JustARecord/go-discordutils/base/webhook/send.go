package webhook

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// SendMessage sends a message to a webhook.
func SendMessage(ctx context.Context, client *discordgo.Session, wh *discordgo.Webhook, params *discordgo.WebhookParams) (*discordgo.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		message, err := client.WebhookExecute(wh.ID, wh.Token, true, params)
		if err != nil {
			return nil, err
		}

		return message, nil
	}
}

// Send sends a message to a webhook (threaded or not).
func Send(ctx context.Context, client *discordgo.Session, wh *discordgo.Webhook, params *discordgo.WebhookParams, thread *discordgo.Channel) (*discordgo.Message, error) {
	if thread == nil {
		return SendMessage(ctx, client, wh, params)
	}

	return SendMessageInThread(ctx, client, wh, params, thread)
}
