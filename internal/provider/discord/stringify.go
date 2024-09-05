package discord

import (
	"sort"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/maps"
)

// Stringify converts a discordgo type to a string.
func Stringify(v interface{}) string {
	var result string
	var ok bool
	switch t := v.(type) {
	// Guild constants
	case discordgo.MessageNotifications:
		result, ok = MessageNotifications[t]
	case discordgo.ExplicitContentFilterLevel:
		result, ok = ExplicitContentFilters[t]
	case discordgo.MfaLevel:
		result, ok = MFALevels[t]
	case discordgo.VerificationLevel:
		result, ok = VerificationLevels[t]
	case discordgo.GuildNSFWLevel:
		result, ok = GuildNSFWLevels[t]
	case discordgo.PremiumTier:
		result, ok = PremiumTiers[t]
	// Channel constants
	case discordgo.ChannelType:
		result, ok = ChannelTypes[t]
	case discordgo.ForumSortOrderType:
		result, ok = ForumSortOrders[t]
	case discordgo.ForumLayout:
		result, ok = ForumLayouts[t]
	// Permission constants
	case discordgo.PermissionOverwriteType:
		result, ok = OverwriteTypes[t]
	// Webhook constants
	case discordgo.WebhookType:
		result, ok = WebhookTypes[t]
	}

	if !ok {
		result = UNIMPLEMENTED
	}

	return result
}

// ListStringify converts a listable discordgo type to a list of strings.
func ListStringify(v interface{}) []string {
	enabled := []string{}

	switch t := v.(type) {
	// Guild constants
	case discordgo.SystemChannelFlag:
		// Fetch the keys and sort them
		keys := maps.Keys(SystemChannelFlags)
		sort.Strings(keys)

		for _, k := range keys {
			v := SystemChannelFlags[k]
			if t&v == v {
				enabled = append(enabled, k)
			}
		}
	// Role constants
	case discordgo.RoleFlags:
		// Fetch the keys and sort them
		keys := maps.Keys(RoleFlags)
		sort.Strings(keys)

		for _, k := range keys {
			v := RoleFlags[k]
			if t&v == v {
				enabled = append(enabled, k)
			}
		}
	case int64:
		// TODO: for now, assume permissions

		// Fetch the keys and sort them
		keys := maps.Keys(Permissions)
		sort.Strings(keys)

		for _, k := range keys {
			v := Permissions[k]
			if t&v == v {
				enabled = append(enabled, k)
			}
		}
	// Channel constants
	case discordgo.ChannelFlags:
		// Fetch the keys and sort them
		keys := maps.Keys(ChannelFlags)
		sort.Strings(keys)

		for _, k := range keys {
			v := ChannelFlags[k]
			if t&v == v {
				enabled = append(enabled, k)
			}
		}
	default:
		enabled = append(enabled, UNIMPLEMENTED)
	}

	return enabled
}
