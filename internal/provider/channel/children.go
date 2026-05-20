package channel

import (
	"context"
	"fmt"
	"strings"

	discordchannel "github.com/JustARecord/go-discordutils/base/channel"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/TheCodedCloud/terraform-provider-discord/internal/provider/common"
)

const (
	ChildrenTypeIDs   = "ids"
	ChildrenTypeNames = "names"
)

var validChildrenTypes = []string{ChildrenTypeIDs, ChildrenTypeNames}

// readChildrenEnabled returns true when read_children is unset (null/unknown) or explicitly true.
func readChildrenEnabled(readChildren types.Bool) bool {
	if readChildren.IsNull() || readChildren.IsUnknown() {
		return true
	}

	return readChildren.ValueBool()
}

// childrenTypeEffective returns "ids" when children_type is unset (null/unknown).
func childrenTypeEffective(childrenType types.String) string {
	if childrenType.IsNull() || childrenType.IsUnknown() {
		return ChildrenTypeIDs
	}

	return strings.ToLower(childrenType.ValueString())
}

func childrenListForCategory(
	ctx context.Context,
	client *discordgo.Session,
	guildID string,
	category *discordgo.Channel,
	readChildren types.Bool,
	childrenType types.String,
) (types.List, diag.Diagnostics) {
	if !readChildrenEnabled(readChildren) {
		return types.ListNull(types.StringType), nil
	}

	children, err := discordchannel.FetchChildren(ctx, client, guildID, category)
	if err != nil {
		return types.ListNull(types.StringType), diag.Diagnostics{
			diag.NewErrorDiagnostic("Failed to get channel children", err.Error()),
		}
	}

	var values []string
	switch childrenTypeEffective(childrenType) {
	case ChildrenTypeIDs:
		values = discordchannel.IDsByPosition(children)
	case ChildrenTypeNames:
		values = discordchannel.NamesByPosition(children)
	default:
		return types.ListNull(types.StringType), diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Invalid children_type",
				fmt.Sprintf(`children_type must be one of: %s`, strings.Join(validChildrenTypes, ", ")),
			),
		}
	}

	childrenList, diags := common.ToListType[string, basetypes.StringType](values)
	if diags.HasError() {
		return types.ListNull(types.StringType), diags
	}

	return childrenList, nil
}
