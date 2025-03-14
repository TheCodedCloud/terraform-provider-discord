package permissions

import (
	"context"
	"fmt"
	"strings"

	discord "github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/thecodedcloud/terraform-provider-discord/internal/provider/common"
)

// UpdateModel updates the plan model with the latest data.
func UpdateModel(ctx context.Context, overwrite *discordgo.PermissionOverwrite, model, state *PermissionsResourceModel) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("Updating %s %s model with model data: %v", resourceMetadataName, resourceMetadataType, model))
	tflog.Info(ctx, fmt.Sprintf("Updating %s %s model with state data: %v", resourceMetadataName, resourceMetadataType, state))

	allow, deny, err := discord.ParseOverwrite(overwrite)
	if err != nil {
		tflog.Error(ctx, "Failed to parse permission overwrite", map[string]interface{}{"overwrite": overwrite, "error": err})
		return diag.Diagnostics{
			diag.NewErrorDiagnostic("failed to parse permission overwrite", err.Error()),
		}
	}

	allowList, diags := common.ToListType[string, basetypes.StringType](allow)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to parse allow list", map[string]interface{}{"allow": allow})
		return diags
	}

	tflog.Info(ctx, fmt.Sprintf("Parsed allow list: %v", allowList))

	denyList, diags := common.ToListType[string, basetypes.StringType](deny)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to parse deny list", map[string]interface{}{"deny": deny})
		return diags
	}

	tflog.Info(ctx, fmt.Sprintf("Parsed deny list: %v", denyList))

	if model == nil {
		model = &PermissionsResourceModel{}
	}

	model.ID = types.StringValue(overwrite.ID)
	model.Type = types.StringValue(strings.ToLower(discord.Stringify(overwrite.Type)))

	if model.Allow.IsNull() || model.Allow.IsUnknown() || len(model.Allow.Elements()) == 0 {
		model.Allow = allowList
	}

	if model.Deny.IsNull() || model.Deny.IsUnknown() || len(model.Deny.Elements()) == 0 {
		model.Deny = denyList
	}

	if state != nil {
		tflog.Info(ctx, "Updating model with state data", map[string]interface{}{"model": model, "state": state})

		if !state.Allow.IsNull() && !state.Allow.IsUnknown() {
			// Revert allow to the plan value
			// This keeps the allow in the same order as the plan
			model.Allow = state.Allow
		}

		if !state.Deny.IsNull() && !state.Deny.IsUnknown() {
			// Revert deny to the plan value
			// This keeps the deny in the same order as the plan
			model.Deny = state.Deny
		}
	}

	if state == nil {
		// If the plan is nil, return early.
		return nil
	}

	// Otherwise, update the model with additional data from the plan.

	// Map the guild data to the state.
	model.GuildID = state.GuildID
	model.ChannelID = state.ChannelID

	return nil
}
