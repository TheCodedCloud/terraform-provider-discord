package utils

import (
	"github.com/bwmarrin/discordgo"
)

const (
	EmbedTitleMaxLen       = 256
	EmbedDescriptionMaxLen = 4096
	EmbedFooterTextMaxLen  = 2048
	EmbedAuthorNameMaxLen  = 256
	EmbedMaxFields         = 25
	EmbedFieldNameMaxLen   = 256
	EmbedFieldValueMaxLen  = 1024
	MessageEmbedsMaxLen    = 6000

	Truncated = "..."
)

var (
	TruncatedLen = len(Truncated)
)

// TruncateEmbed ensures that the embed character limits are not exceeded.
// The limits are as follows:
// - Title: 256 characters
// - Description: 4096 characters
// - Fields: 25 fields
// - Field name: 256 characters
// - Field value: 1024 characters
// - Footer text: 2048 characters
// - Author name: 256 characters
func TruncateEmbed(embed *discordgo.MessageEmbed) {
	// Make sure truncated fields end with "..." and are still within the character limit
	if len(embed.Title) > EmbedTitleMaxLen {
		embed.Title = embed.Title[:(EmbedTitleMaxLen-TruncatedLen)] + Truncated
	}

	if len(embed.Description) > EmbedDescriptionMaxLen {
		embed.Description = embed.Description[:(EmbedDescriptionMaxLen-TruncatedLen)] + Truncated
	}

	if embed.Footer != nil && len(embed.Footer.Text) > EmbedFooterTextMaxLen {
		embed.Footer.Text = embed.Footer.Text[:(EmbedFooterTextMaxLen-TruncatedLen)] + Truncated
	}

	if embed.Author != nil && len(embed.Author.Name) > EmbedAuthorNameMaxLen {
		embed.Author.Name = embed.Author.Name[:(EmbedAuthorNameMaxLen-TruncatedLen)] + Truncated
	}

	if len(embed.Fields) > EmbedMaxFields {
		embed.Fields = embed.Fields[:EmbedMaxFields]
	}

	for _, field := range embed.Fields {
		if len(field.Name) > EmbedFieldNameMaxLen {
			field.Name = field.Name[:(EmbedFieldNameMaxLen-TruncatedLen)] + Truncated
		}

		if len(field.Value) > EmbedFieldValueMaxLen {
			field.Value = field.Value[:(EmbedFieldValueMaxLen-TruncatedLen)] + Truncated
		}
	}
}
