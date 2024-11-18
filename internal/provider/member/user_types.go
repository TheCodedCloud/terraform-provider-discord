package member

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// User defines the schema for the user.
type User struct {
	// The user's ID.
	ID types.String `tfsdk:"id"`

	// The user's username, not unique across the platform.
	Username types.String `tfsdk:"username"`

	// The user's Discord-tag.
	Discriminator types.String `tfsdk:"discriminator"`

	// The user's display name, if set. For bots, this is the application name.
	GlobalName types.String `tfsdk:"global_name"`

	// Whether the user belongs to an OAuth2 application.
	Bot types.Bool `tfsdk:"bot"`

	// Whether the user is an Official Discord System user (part of the urgent message system).
	System types.Bool `tfsdk:"system"`

	// Whether the user has two factor enabled on their account.
	MFAEnabled types.Bool `tfsdk:"mfa_enabled"`

	// The user's banner hash.
	Banner types.String `tfsdk:"banner"`

	// The user's banner color if they have one.
	AccentColor types.String `tfsdk:"accent_color"`

	// The user's chosen language option.
	Locale types.String `tfsdk:"locale"`

	// Whether the email on this account has been verified.
	Verified types.Bool `tfsdk:"verified"`

	// The user's email.
	Email types.String `tfsdk:"email"`

	// The flags on a user's account.
	Flags types.List `tfsdk:"flags"`

	// The type of Nitro subscription on a user's account.
	PremiumType types.String `tfsdk:"premium_type"`

	// The public flags on a user's account.
	PublicFlags types.List `tfsdk:"public_flags"`

	// Data for the user's avatar decoration.
	// AvatarDecorationData types.String `tfsdk:"avatar_decoration_data"`
	Avatar types.String `tfsdk:"avatar"`
}

var UserSchema = schema.SingleNestedAttribute{
	Description: "The user object.",
	Computed:    true,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"username": schema.StringAttribute{
			Computed: true,
		},
		"discriminator": schema.StringAttribute{
			Computed: true,
		},
		"global_name": schema.StringAttribute{
			Computed: true,
		},
		"bot": schema.BoolAttribute{
			Computed: true,
		},
		"system": schema.BoolAttribute{
			Computed: true,
		},
		"mfa_enabled": schema.BoolAttribute{
			Computed: true,
		},
		"banner": schema.StringAttribute{
			Computed: true,
		},
		"accent_color": schema.StringAttribute{
			Computed: true,
		},
		"locale": schema.StringAttribute{
			Computed: true,
		},
		"verified": schema.BoolAttribute{
			Computed: true,
		},
		"email": schema.StringAttribute{
			Computed: true,
		},
		"flags": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"premium_type": schema.StringAttribute{
			Computed: true,
		},
		"public_flags": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"avatar": schema.StringAttribute{
			Computed: true,
		},
	},
}
