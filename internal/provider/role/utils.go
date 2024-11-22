package role

import (
	discord "github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
)

func setupParams(model *RoleResourceModel, permissions []string) *discordgo.RoleParams {
	name := model.Name.ValueString()

	roleParams := &discordgo.RoleParams{
		Name: name,
	}

	if len(permissions) > 0 {
		permissionsSum := discord.CalcPermissions(permissions)
		roleParams.Permissions = &permissionsSum
	}

	// Set optional parameters
	if !model.Color.IsNull() {
		color := common.IntHex(model.Color.ValueString())
		roleParams.Color = &color
	}

	if !model.Hoist.IsNull() {
		roleParams.Hoist = model.Hoist.ValueBoolPointer()
	}

	if !model.Mentionable.IsNull() {
		roleParams.Mentionable = model.Mentionable.ValueBoolPointer()
	}

	if !model.UnicodeEmoji.IsNull() {
		roleParams.UnicodeEmoji = model.UnicodeEmoji.ValueStringPointer()
	}

	if !model.Icon.IsNull() {
		roleParams.Icon = model.Icon.ValueStringPointer()
	}

	return roleParams
}

// UpdateModel updates the role resource model with the provided role.
func UpdateModel(role *discordgo.Role, model, state *RoleResourceModel) diag.Diagnostics {
	permissions := discord.ListStringify(role.Permissions)
	flags := discord.ListStringify(role.Flags)

	permissionsList, diags := common.ToListType[string, basetypes.StringType](permissions)
	if diags.HasError() {
		return diags
	}

	flagsList, diags := common.ToListType[string, basetypes.StringType](flags)
	if diags.HasError() {
		return diags
	}

	if model == nil {
		model = &RoleResourceModel{}
	}

	model.ID = types.StringValue(role.ID)
	model.Name = types.StringValue(role.Name)
	model.Color = types.StringValue(common.StrHex(role.Color))
	model.Hoist = types.BoolValue(role.Hoist)
	model.Position = types.Int32Value(int32(role.Position))
	model.Managed = types.BoolValue(role.Managed)
	model.Mentionable = types.BoolValue(role.Mentionable)
	model.Icon = types.StringValue(role.Icon)
	model.UnicodeEmoji = types.StringValue(role.UnicodeEmoji)
	model.Flags = flagsList

	if model.Permissions.IsNull() {
		model.Permissions = permissionsList
	}

	if state != nil && !state.Permissions.IsNull() {
		// Revert permissions to the plan value
		// This keeps the permissions in the same order as the plan
		model.Permissions = state.Permissions
	}

	if state == nil {
		// If the plan is nil, return early.
		return nil
	}

	// Otherwise, update the model with additional data from the plan.

	// Map the guild data to the state.
	model.GuildID = state.GuildID

	return nil
}
