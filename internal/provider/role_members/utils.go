package role_members

import (
	"github.com/JustARecord/go-discordutils/base/member"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
)

// UpdateModel updates the role resource model with the provided role.
func UpdateModel(role *discordgo.Role, members []*discordgo.Member, model, state *RoleMembersResourceModel) diag.Diagnostics {
	memberNames := member.Names(members)
	membersList, diags := common.ToListType[string, basetypes.StringType](memberNames)
	if diags.HasError() {
		return diags
	}

	if model == nil {
		model = &RoleMembersResourceModel{}
	}

	model.RoleID = types.StringValue(role.ID)
	model.Role = types.StringValue(role.Name)

	if model.Members.IsNull() {
		model.Members = membersList
	}

	// if state != nil && !state.Members.IsNull() {
	// 	// Revert members to the plan value
	// 	// This keeps the members in the same order as the plan
	// 	model.Members = state.Members
	// }

	if state == nil {
		// If the plan is nil, return early.
		return nil
	}

	// Otherwise, update the model with additional data from the plan.

	// Map the guild data to the state.
	model.GuildID = state.GuildID

	return nil
}
