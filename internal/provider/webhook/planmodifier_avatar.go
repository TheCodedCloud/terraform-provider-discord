package webhook

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// avatarUnknownIfImageData returns a plan modifier that sets the planned avatar value to
// unknown when the config value is image data (data URL or base64). Discord accepts image
// data but returns a hash; so the value is effectively computed. Setting plan to unknown
// avoids "Provider produced inconsistent result" when we store the hash in state.
func avatarUnknownIfImageData() planmodifier.String {
	return avatarUnknownIfImageDataModifier{}
}

type avatarUnknownIfImageDataModifier struct{}

func (avatarUnknownIfImageDataModifier) Description(context.Context) string {
	return "If the config value is image data (e.g. data URL or base64), the planned value is set to unknown so the provider can store the hash Discord returns."
}

func (avatarUnknownIfImageDataModifier) MarkdownDescription(context.Context) string {
	return "If the config value is image data (e.g. data URL or base64), the planned value is set to unknown so the provider can store the hash Discord returns."
}

func (m avatarUnknownIfImageDataModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	v := req.ConfigValue.ValueString()
	if !isAvatarImageData(v) {
		return
	}
	resp.PlanValue = types.StringUnknown()
}
