package channel

import (
	"sort"

	"github.com/bwmarrin/discordgo"
)

// IDs returns channel IDs sorted alphabetically for stable comparisons.
// Use IDsByPosition when reflecting Discord sidebar order.
func IDs(channels []*discordgo.Channel) []string {
	ids := make([]string, len(channels))
	for i, c := range channels {
		ids[i] = c.ID
	}

	sort.Strings(ids)

	return ids
}

// Names returns channel names sorted alphabetically for stable comparisons.
// Use NamesByPosition when reflecting Discord sidebar order.
func Names(channels []*discordgo.Channel) []string {
	names := make([]string, len(channels))
	for i, c := range channels {
		names[i] = c.Name
	}

	sort.Strings(names)

	return names
}

// SortByPosition returns a copy of channels sorted by Position ascending (Discord sidebar order).
func SortByPosition(channels []*discordgo.Channel) []*discordgo.Channel {
	if len(channels) == 0 {
		return nil
	}

	sorted := make([]*discordgo.Channel, len(channels))
	copy(sorted, channels)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Position < sorted[j].Position
	})

	return sorted
}

// NamesByPosition returns channel names in Discord sidebar order (by Position ascending).
func NamesByPosition(channels []*discordgo.Channel) []string {
	sorted := SortByPosition(channels)
	names := make([]string, len(sorted))
	for i, c := range sorted {
		names[i] = c.Name
	}

	return names
}

// IDsByPosition returns channel IDs in Discord sidebar order (by Position ascending).
func IDsByPosition(channels []*discordgo.Channel) []string {
	sorted := SortByPosition(channels)
	ids := make([]string, len(sorted))
	for i, c := range sorted {
		ids[i] = c.ID
	}

	return ids
}
