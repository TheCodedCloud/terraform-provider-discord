package channel

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ForumTag struct {
	// The ID of the tag.
	ID types.String `tfsdk:"id"`

	// The name of the tag (0-20 characters)
	Name types.String `tfsdk:"name"`

	// Whether this tag can only be added to or removed from threads by a member with the MANAGE_THREADS permission
	Moderated types.Bool `tfsdk:"moderated"`

	// The ID of a guild's custom emoji.
	EmojiID types.String `tfsdk:"emoji_id"`

	// The unicode character of the emoji.
	EmojiName types.String `tfsdk:"emoji_name"`

	// Note: only one of either emoji_id or emoji_name can be set
}

var ForumTagSchema = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"moderated": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
		"emoji_id": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"emoji_name": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
	},
}

type ForumDefaultReaction struct {
	// The ID of a guild's custom emoji.
	EmojiID types.String `tfsdk:"emoji_id"`

	// The unicode character of the emoji.
	EmojiName types.String `tfsdk:"emoji_name"`

	// Note: only one of either emoji_id or emoji_name can be set
}

type ForumSortOrderType struct {
}
