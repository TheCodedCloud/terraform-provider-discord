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
// We store exactly what we sent so the plan (config) matches state and Terraform does not report
// "planned value does not match config value" or "inconsistent result". When the user sends image
// data (base64/data URL), we keep it in state; when they send "" we store "".
func avatarStateAfterApply(sentValue, _ string) types.String {
	return types.StringValue(sentValue)
}

// isAvatarImageData returns true if the value looks like image payload (data URL or raw base64).
// When state contains image data, Read keeps it (we do not overwrite with the hash from the API)
// so that the next plan still matches state.
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
