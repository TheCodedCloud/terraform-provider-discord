package webhook

import (
	"context"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

// CreateThreadWithMessages creates a thread with messages.
func CreateThreadWithMessages(ctx context.Context, client *discordgo.Session, wh *discordgo.Webhook, channel *discordgo.Channel, name string, messages []*discordgo.WebhookParams) (*discordgo.Channel, error) {
	// 1. Send the message to the channel
	initialMessage := messages[0]
	messages = messages[1:]

	message, err := SendMessage(ctx, client, wh, initialMessage)
	if err != nil {
		return nil, err
	}

	// 2. Create the thread with the message
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		thread, err := client.MessageThreadStartComplex(channel.ID, message.ID, &discordgo.ThreadStart{
			Name: name,
		})
		if err != nil {
			return nil, err
		}

		// 3. Send the remaining messages to the thread
		for _, message := range messages {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if _, err := SendMessageInThread(ctx, client, wh, message, thread); err != nil {
					return nil, err
				}
			}
		}

		return thread, nil
	}
}

// SendMessageInThread sends a message to a thread.
func SendMessageInThread(ctx context.Context, client *discordgo.Session, wh *discordgo.Webhook, params *discordgo.WebhookParams, thread *discordgo.Channel) (*discordgo.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		message, err := client.WebhookThreadExecute(wh.ID, wh.Token, true, thread.ID, params)
		if err != nil {
			return nil, err
		}

		return message, nil
	}
}

// EditThreadMessage edits a thread message.
func EditThreadMessage(ctx context.Context, client *discordgo.Session, webhook *discordgo.Webhook, thread *discordgo.Channel, messageID string, params *discordgo.WebhookEdit) (*discordgo.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		uri := discordgo.EndpointWebhookMessage(webhook.ID, webhook.Token, messageID)

		// store the thread ID in query params
		uri += "?thread_id=" + thread.ID

		var response []byte
		var err error
		if len(params.Files) > 0 {
			_, body, err := discordgo.MultipartBodyWithJSON(params, params.Files)
			if err != nil {
				return nil, err
			}

			response, err = client.Request("PATCH", uri, body)
			if err != nil {
				return nil, err
			}
		} else {
			response, err = client.RequestWithBucketID("PATCH", uri, params, discordgo.EndpointWebhookToken("", ""))
			if err != nil {
				return nil, err
			}
		}

		var msg *discordgo.Message
		if err = json.Unmarshal(response, &msg); err != nil {
			return nil, err
		}

		return msg, nil
	}
}
