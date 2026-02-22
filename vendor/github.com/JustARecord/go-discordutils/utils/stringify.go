package utils

import (
	"sort"

	"github.com/JustARecord/go-discordutils/base/common"
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

// Stringify converts a discordgo type to a string.
func Stringify(v interface{}) string {
	var result string
	var ok bool
	switch t := v.(type) {
	// Guild constants
	case discordgo.MessageNotifications:
		result, ok = common.MessageNotifications[t]
	case discordgo.ExplicitContentFilterLevel:
		result, ok = common.ExplicitContentFilters[t]
	case discordgo.MfaLevel:
		result, ok = common.MFALevels[t]
	case discordgo.VerificationLevel:
		result, ok = common.VerificationLevels[t]
	case discordgo.GuildNSFWLevel:
		result, ok = common.GuildNSFWLevels[t]
	case discordgo.PremiumTier:
		result, ok = common.PremiumTiers[t]
	// Channel constants
	case discordgo.ChannelType:
		result, ok = common.ChannelTypes[t]
	case discordgo.ForumSortOrderType:
		result, ok = common.ForumSortOrders[t]
	case discordgo.ForumLayout:
		result, ok = common.ForumLayouts[t]
	// Permission constants
	case discordgo.PermissionOverwriteType:
		result, ok = common.OverwriteTypes[t]
	// Webhook constants
	case discordgo.WebhookType:
		result, ok = common.WebhookTypes[t]
	// User constants
	case discordgo.UserPremiumType:
		result, ok = common.UserPremiumTypes[t]
	}

	if !ok {
		result = common.UNIMPLEMENTED
	}

	return result
}

// KeyStringify fetches the key of a discordgo type map that matches the value.
func KeyStringify[K comparable](kType interface{}, v interface{}) (interface{}, bool) {
	var result interface{}
	var ok bool

	var vString string
	switch t := v.(type) {
	case string:
		vString = t
	default:
		vString = Stringify(v)
	}

	// Switch on the type of the map
	switch kType.(type) {
	// Guild constants
	case discordgo.MessageNotifications:
		result, ok = FetchKeyByValue(common.MessageNotifications, vString)
	case discordgo.ExplicitContentFilterLevel:
		result, ok = FetchKeyByValue(common.ExplicitContentFilters, vString)
	case discordgo.MfaLevel:
		result, ok = FetchKeyByValue(common.MFALevels, vString)
	case discordgo.VerificationLevel:
		result, ok = FetchKeyByValue(common.VerificationLevels, vString)
	case discordgo.GuildNSFWLevel:
		result, ok = FetchKeyByValue(common.GuildNSFWLevels, vString)
	case discordgo.PremiumTier:
		result, ok = FetchKeyByValue(common.PremiumTiers, vString)
	// Channel constants
	case discordgo.ChannelType:
		result, ok = FetchKeyByValue(common.ChannelTypes, vString)
	case discordgo.ForumSortOrderType:
		result, ok = FetchKeyByValue(common.ForumSortOrders, vString)
	case discordgo.ForumLayout:
		result, ok = FetchKeyByValue(common.ForumLayouts, vString)
	// Permission constants
	case discordgo.PermissionOverwriteType:
		result, ok = FetchKeyByValue(common.OverwriteTypes, vString)
	// Webhook constants
	case discordgo.WebhookType:
		result, ok = FetchKeyByValue(common.WebhookTypes, vString)
	// User constants
	case discordgo.UserPremiumType:
		result, ok = FetchKeyByValue(common.UserPremiumTypes, vString)
	}

	if !ok {
		result = common.UNIMPLEMENTED
	}

	return result, ok
}

// ListStringify converts a listable discordgo type to a list of strings.
func ListStringify(v interface{}) []string {
	enabled := []string{}

	switch t := v.(type) {
	// Guild constants
	case discordgo.SystemChannelFlag:
		// Fetch the keys and sort them
		keys := maps.Keys(common.SystemChannelFlags)
		sort.Strings(keys)

		enabled = lo.Filter(keys, func(k string, _ int) bool {
			v := common.SystemChannelFlags[k]
			return t&v == v
		})
	// Role constants
	case discordgo.RoleFlags:
		// Fetch the keys and sort them
		keys := maps.Keys(common.RoleFlags)
		sort.Strings(keys)

		enabled = lo.Filter(keys, func(k string, _ int) bool {
			v := common.RoleFlags[k]
			return t&v == v
		})
	case int64:
		// TODO: for now, assume permissions

		// Fetch the keys and sort them
		keys := maps.Keys(common.Permissions)
		sort.Strings(keys)

		enabled = lo.Filter(keys, func(k string, _ int) bool {
			v := common.Permissions[k]
			return t&v == v
		})
	// Channel constants
	case discordgo.ChannelFlags:
		// Fetch the keys and sort them
		keys := maps.Keys(common.ChannelFlags)
		sort.Strings(keys)

		enabled = lo.Filter(keys, func(k string, _ int) bool {
			v := common.ChannelFlags[k]
			return t&v == v
		})
	// Member constants
	case discordgo.MemberFlags:
		// Fetch the keys and sort them
		keys := maps.Keys(common.MemberFlags)
		sort.Strings(keys)

		enabled = lo.Filter(keys, func(k string, _ int) bool {
			v := common.MemberFlags[k]
			return t&v == v
		})
	case discordgo.UserFlags:
		// TODO: for now, assume user flags

		// Fetch the keys and sort them
		keys := maps.Keys(common.UserFlags)
		sort.Strings(keys)

		enabled = lo.Filter(keys, func(k string, _ int) bool {
			v := common.UserFlags[k]
			return t&v == v
		})
	default:
		enabled = append(enabled, common.UNIMPLEMENTED)
	}

	return enabled
}
