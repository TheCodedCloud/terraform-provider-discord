package channel

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

// EditMessageByID edits a channel message.
func EditMessageByID(ctx context.Context, client *discordgo.Session, channelID, messageID string, params *discordgo.MessageEdit) (*discordgo.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Make sure the channel ID and message ID are set
		params.Channel = channelID
		params.ID = messageID

		return client.ChannelMessageEditComplex(params)
	}
}

// EditMessage edits a channel message.
func EditMessage(ctx context.Context, client *discordgo.Session, channel *discordgo.Channel, messageID string, params *discordgo.MessageEdit) (*discordgo.Message, error) {
	return EditMessageByID(ctx, client, channel.ID, messageID, params)
}

// FetchMessageByID fetches a channel message.
func FetchMessageByID(ctx context.Context, client *discordgo.Session, channel_id, message_id string) (*discordgo.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return client.ChannelMessage(channel_id, message_id)
	}
}

// FetchMessage fetches a channel message by channel.
func FetchMessage(ctx context.Context, client *discordgo.Session, channel *discordgo.Channel, message_id string) (*discordgo.Message, error) {
	return FetchMessageByID(ctx, client, channel.ID, message_id)
}

// FetchMessagesByID fetches messages in a channel.
func FetchMessagesByID(ctx context.Context, client *discordgo.Session, channel_id string, limit int) ([]*discordgo.Message, error) {
	if limit <= 100 {
		return client.ChannelMessages(channel_id, limit, "", "", "")
	}

	// Fetch messages in chunks of 100, until we reach the first message
	var messages []*discordgo.Message

messagesLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			lastMessage := messages[len(messages)-1]
			chunk, err := client.ChannelMessages(channel_id, 100, lastMessage.ID, "", "")
			if err != nil {
				return nil, err
			}

			messages = append(messages, chunk...)
			if len(chunk) < 100 {
				break messagesLoop
			}
		}
	}

	return messages, nil
}

// FetchMessages fetches messages in a channel.
func FetchMessages(ctx context.Context, client *discordgo.Session, channel *discordgo.Channel, limit int) ([]*discordgo.Message, error) {
	return FetchMessagesByID(ctx, client, channel.ID, limit)
}

// FetchMessageIDs fetches message IDs in a channel.
func FetchMessageIDs(ctx context.Context, client *discordgo.Session, channel_id string, limit int) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		messages, err := FetchMessagesByID(ctx, client, channel_id, limit)
		if err != nil {
			return nil, err
		}

		ids := lo.Map(messages, func(m *discordgo.Message, _ int) string {
			return m.ID
		})

		return ids, nil
	}
}

// DeleteMessageByID deletes a channel message.
func DeleteMessageByID(ctx context.Context, client *discordgo.Session, channel_id, message_id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return client.ChannelMessageDelete(channel_id, message_id)
	}
}

// DeleteMessage deletes a channel message.
func DeleteMessage(ctx context.Context, client *discordgo.Session, channel *discordgo.Channel, message_id string) error {
	return DeleteMessageByID(ctx, client, channel.ID, message_id)
}

// DeleteMessages deletes messages in a channel (by message IDs).
func DeleteMessages(ctx context.Context, client *discordgo.Session, channel_id string, messages []string) error {
	if len(messages) > 100 {
		return fmt.Errorf("cannot bulk delete more than 100 messages")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return client.ChannelMessagesBulkDelete(channel_id, messages)
	}
}

// Empty empties a channel.
func Empty(ctx context.Context, client *discordgo.Session, channel_id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Fetch first 100 messages
		messages, err := FetchMessageIDs(ctx, client, channel_id, 100)
		if err != nil {
			return err
		}

		for len(messages) > 0 {
			// Bulk delete messages
			if err := DeleteMessages(ctx, client, channel_id, messages); err != nil {
				return err
			}

			// Fetch next 100 messages
			messages, err = FetchMessageIDs(ctx, client, channel_id, 100)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
