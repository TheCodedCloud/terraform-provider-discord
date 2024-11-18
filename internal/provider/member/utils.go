package member

import (
	"context"

	discord "github.com/JustARecord/go-discordutils/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/justarecord/terraform-provider-discord/internal/provider/common"
)

// ToUser converts a discordgo.User to a User.
func ToUser(ctx context.Context, user *discordgo.User, guildID string) (*User, diag.Diagnostics) {
	accentColor := common.StrHex(user.AccentColor)

	flags := discord.ListStringify(user.Flags)
	flagsList, diags := common.ToListType[string, basetypes.StringType](flags)
	if diags.HasError() {
		return nil, diags
	}

	premiumType := discord.Stringify(user.PremiumType)

	publicFlags := discord.ListStringify(user.PublicFlags)
	publicFlagsList, diags := common.ToListType[string, basetypes.StringType](publicFlags)
	if diags.HasError() {
		return nil, diags
	}

	result := &User{
		ID:            types.StringValue(user.ID),
		Email:         types.StringValue(user.Email),
		Username:      types.StringValue(user.Username),
		Avatar:        types.StringValue(user.Avatar),
		Locale:        types.StringValue(user.Locale),
		Discriminator: types.StringValue(user.Discriminator),
		GlobalName:    types.StringValue(user.GlobalName),
		Verified:      types.BoolValue(user.Verified),
		MFAEnabled:    types.BoolValue(user.MFAEnabled),
		Banner:        types.StringValue(user.Banner),
		AccentColor:   types.StringValue(accentColor),
		Bot:           types.BoolValue(user.Bot),
		PublicFlags:   publicFlagsList,
		PremiumType:   types.StringValue(premiumType),
		System:        types.BoolValue(user.System),
		Flags:         flagsList,
	}

	return result, nil
}
