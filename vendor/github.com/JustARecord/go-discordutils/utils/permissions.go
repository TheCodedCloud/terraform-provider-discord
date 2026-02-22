package utils

import (
	"sort"

	"github.com/JustARecord/go-discordutils/base/common"
	"github.com/bwmarrin/discordgo"
)

// HasPermission checks if the provided permissions have the required permission.
func HasPermission(perms int64, required int64) bool {
	return perms&required == required
}

// CalcPermissions calculates the permissions from a list of permissions.
func CalcPermissions(permissions []string) int64 {
	var perms int64
	for _, perm := range permissions {
		perms |= common.Permissions[perm]
	}
	return perms
}

// ParseOverwrite parses the provided permission overwrite into allowed and denied permissions lists.
func ParseOverwrite(overwrite *discordgo.PermissionOverwrite) ([]string, []string, error) {
	allowed := ListStringify(overwrite.Allow)
	denied := ListStringify(overwrite.Deny)

	sort.Strings(allowed)
	sort.Strings(denied)

	return allowed, denied, nil
}
