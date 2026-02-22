package role

import (
	"sort"

	"github.com/bwmarrin/discordgo"
)

// Names returns a list of role names.
func Names(roles []*discordgo.Role) []string {
	names := []string{}
	for _, r := range roles {
		names = append(names, r.Name)
	}

	// Sort the names
	sort.Strings(names)

	return names
}

// IDs returns a list of role IDs.
func IDs(roles []*discordgo.Role) []string {
	ids := []string{}
	for _, r := range roles {
		ids = append(ids, r.ID)
	}

	// Sort the IDs
	sort.Strings(ids)

	return ids
}

// SortByPosition sorts the roles by their position.
func SortByPosition(roles []*discordgo.Role) {
	if len(roles) == 0 {
		return
	}

	// Bubble sort the roles by position (descending)
	for i := 0; i < len(roles); i++ {
		for j := i + 1; j < len(roles); j++ {
			if roles[i].Position < roles[j].Position {
				roles[i], roles[j] = roles[j], roles[i]
			}
		}
	}
}
