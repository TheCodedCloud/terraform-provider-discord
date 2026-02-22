package webhook

import "github.com/bwmarrin/discordgo"

// ToWebhookEdit converts a message edit to a webhook edit.
func ToWebhookEdit(params *discordgo.MessageEdit) *discordgo.WebhookEdit {
	return &discordgo.WebhookEdit{
		Content:         params.Content,
		Components:      params.Components,
		Embeds:          params.Embeds,
		Files:           params.Files,
		Attachments:     params.Attachments,
		AllowedMentions: params.AllowedMentions,
	}
}
