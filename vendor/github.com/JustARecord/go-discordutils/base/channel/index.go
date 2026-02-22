package channel

import (
	"sort"

	"github.com/bwmarrin/discordgo"
)

// IDs returns a list of channel IDs.
func IDs(channels []*discordgo.Channel) []string {
	ids := make([]string, len(channels))
	for i, c := range channels {
		ids[i] = c.ID
	}

	// Sort the IDs
	sort.Strings(ids)

	return ids
}

// Names returns a list of channel names.
func Names(channels []*discordgo.Channel) []string {
	names := make([]string, len(channels))
	for i, c := range channels {
		names[i] = c.Name
	}

	// Sort the names
	sort.Strings(names)

	return names
}
