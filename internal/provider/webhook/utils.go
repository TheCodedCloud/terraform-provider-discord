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
