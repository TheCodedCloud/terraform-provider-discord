package webhook

import (
	"context"
	"fmt"
	"strings"

	discord "github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// UpdateModel updates the plan model with the latest data.
func UpdateModel(ctx context.Context, webhook *discordgo.Webhook, model, state *WebhookResourceModel) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("Updating %s %s model with model data: %v", resourceMetadataName, resourceMetadataType, model))
	tflog.Info(ctx, fmt.Sprintf("Updating %s %s model with state data: %v", resourceMetadataName, resourceMetadataType, state))

	if model == nil {
		model = &WebhookResourceModel{}
	}

	model.ID = types.StringValue(webhook.ID)
	model.Type = types.StringValue(strings.ToLower(discord.Stringify(webhook.Type)))
	model.Name = types.StringValue(webhook.Name)
	model.Token = types.StringValue(webhook.Token)
	model.ApplicationID = types.StringValue(webhook.ApplicationID)
	model.GuildID = types.StringValue(webhook.GuildID)
	model.ChannelID = types.StringValue(webhook.ChannelID)

	if state == nil {
		return nil
	}

	// Set the Avatar from the state if it is not empty.
	model.Avatar = state.Avatar

	return nil
}

// avatarStateAfterApply returns the value to store in state for avatar after a successful Create or Update.
// We keep state consistent with what we sent: when we sent no custom avatar (empty), store "";
// when we sent custom avatar data, store the hash Discord returned. We do not rely on Discord's
// response format for "no avatar" (e.g. null vs omit vs default hash) so apply always succeeds
// and state stays consistent.
func avatarStateAfterApply(sentValue, apiAvatarHash string) types.String {
	if sentValue == "" {
		return types.StringValue("")
	}
	return types.StringValue(apiAvatarHash)
}

// isAvatarImageData returns true if the value looks like image payload (data URL or raw base64)
// that Discord will convert to a hash. We use this in plan modification so the plan shows
// "(known after apply)" instead of the raw value, avoiding inconsistent result errors.
func isAvatarImageData(s string) bool {
	if s == "" {
		return false
	}
	if strings.HasPrefix(s, "data:") {
		return true
	}
	// Discord accepts raw base64; hashes are short hex-like strings (e.g. 32 chars).
	if len(s) > 64 {
		return true
	}
	return false
}
