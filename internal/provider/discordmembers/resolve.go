package discordmembers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/JustARecord/go-discordutils/base/role"
	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// searchLimit is the max results requested per GuildMembersSearch call. Discord
// caps this at 1000; there's no reason to request fewer, since the exact-match
// filter below runs locally regardless of how many prefix matches come back.
const searchLimit = 1000

// ResolveMembersByUsername resolves each declared username to a discordgo.Member
// via the Search Guild Members endpoint (GuildMembersSearch), one call per
// username, followed by an exact-match filter — Search matches by *prefix* on
// username/nickname, so a query of "al" would also match "alice"/"albert"
// without this filter. This replaces go-discordutils's guild.FetchMembersByName,
// which resolves names by paginating the gated List Guild Members endpoint and
// filtering the full guild roster in memory.
//
// A username with no exact match produces a warning diagnostic — unlike the
// function it replaces, which silently drops unmatched names via lo.Filter —
// so a typo'd entry in a role's members file is visible rather than quietly
// ignored. This is a deliberate behavior improvement, not just parity.
func ResolveMembersByUsername(ctx context.Context, client DiscordMemberClient, guildID string, usernames []string) ([]*discordgo.Member, diag.Diagnostics) {
	var diags diag.Diagnostics
	resolved := make([]*discordgo.Member, 0, len(usernames))

	for _, username := range usernames {
		select {
		case <-ctx.Done():
			diags.AddError("Context cancelled while resolving members", ctx.Err().Error())
			return nil, diags
		default:
		}

		results, err := client.GuildMembersSearch(guildID, username, searchLimit)
		if err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to search guild members for %q", username),
				err.Error(),
			)
			continue
		}

		found := false
		for _, m := range results {
			if m != nil && m.User != nil && m.User.Username == username {
				resolved = append(resolved, m)
				found = true
				break
			}
		}

		if !found {
			diags.AddWarning(
				"Declared member not found",
				fmt.Sprintf("No guild member with username %q was found. It will not be added to (or will be removed from) this role. Check for a typo in the members list, or confirm the user is still a member of the guild.", username),
			)
		}
	}

	return resolved, diags
}

// FetchDeclaredHolders answers "of these declared members, which currently hold
// role roleID" using only the ungated per-member GET (GuildMember) — one call
// per declared member. This replaces go-discordutils's role.FetchMembers, which
// answers the equivalent guild-wide question ("who in the whole guild holds
// this role") by paginating the gated List Guild Members endpoint and filtering
// client-side. Discord has no endpoint that lists members by role at any
// privilege level, so restricting the question to a declared set is the only
// way to answer a bounded version of it without GUILD_MEMBERS.
//
// This is the mechanism that actually removes this resource's dependency on the
// privileged intent — it can only ever report on drift among members already
// declared here. A member granted this role by hand outside Terraform, and
// never declared in this resource's config, is invisible to this check and will
// not be detected or removed. That narrowing of drift-detection scope is
// deliberate and accepted, not an oversight — see the package/resource docs.
//
// A declared member who has left the guild since the last apply (a 404 from
// GuildMember) is dropped from the result and reported as a warning, not
// silently dropped and not a hard failure of the whole call. Any other error
// (rate limit, auth, 5xx) is a hard failure: it must not be mistaken for "the
// member left," which would otherwise plan a spurious removal on a transient
// hiccup.
func FetchDeclaredHolders(ctx context.Context, client DiscordMemberClient, guildID, roleID string, declared []*discordgo.Member) ([]*discordgo.Member, diag.Diagnostics) {
	var diags diag.Diagnostics
	held := make([]*discordgo.Member, 0, len(declared))

	for _, m := range declared {
		if m == nil || m.User == nil {
			continue
		}

		select {
		case <-ctx.Done():
			diags.AddError("Context cancelled while checking role membership", ctx.Err().Error())
			return nil, diags
		default:
		}

		fresh, err := client.GuildMember(guildID, m.User.ID)
		if err != nil {
			var restErr *discordgo.RESTError
			if errors.As(err, &restErr) && restErr.Response != nil && restErr.Response.StatusCode == http.StatusNotFound {
				diags.AddWarning(
					"Declared member left the guild",
					fmt.Sprintf("%q (id %s) is no longer a member of this guild. Removing them from this role's tracked membership.", m.User.Username, m.User.ID),
				)
				continue
			}

			diags.AddError(
				fmt.Sprintf("Failed to check role membership for %q", m.User.Username),
				err.Error(),
			)
			return nil, diags
		}

		if fresh != nil && slices.Contains(fresh.Roles, roleID) {
			held = append(held, fresh)
		}
	}

	return held, diags
}

// SyncDeclaredMembers reconciles a role's membership against a desired set of
// members, using only ungated per-member Discord calls. It checks current role
// membership across the union of `previous` (the members this resource
// declared as of the last apply — empty on Create) and `desired` (what the
// current plan declares), then adds/removes members so that membership matches
// `desired` exactly among that checked set.
//
// Checking previous ∪ desired — rather than "everyone in the guild who holds
// this role" — is what lets this stay within FetchDeclaredHolders's declared-set
// model while still correctly handling removal: a member deleted from the
// members list between applies is in `previous` but not `desired`, so they're
// included in the check, found to still hold the role, and removed. A member
// hand-granted the role outside Terraform was never declared here in the first
// place, so they're outside both sets and untouched — the same accepted
// narrowing FetchDeclaredHolders documents.
//
// Unlike go-discordutils's role.SetMembers (which this replaces), there is no
// third fetch to "refresh" the result after writing: once toAdd/toRemove have
// been applied to the checked set, the resulting membership is exactly
// `desired` by construction.
func SyncDeclaredMembers(ctx context.Context, client *discordgo.Session, guild *discordgo.Guild, targetRole *discordgo.Role, previous, desired []*discordgo.Member) ([]*discordgo.Member, diag.Diagnostics) {
	var diags diag.Diagnostics

	tracked := unionMembersByID(previous, desired)

	held, holderDiags := FetchDeclaredHolders(ctx, NewSessionMemberClient(client), guild.ID, targetRole.ID, tracked)
	diags.Append(holderDiags...)
	if diags.HasError() {
		return nil, diags
	}

	toAdd, toRemove := partitionForSync(desired, held)

	// go-discordutils's role.AddMembers/RemoveMembers are already ungated — they
	// only call member.AddRole/RemoveRole per member (GuildMemberRoleAdd/Remove)
	// — so they're reused unmodified rather than reimplemented here.
	if len(toRemove) > 0 {
		if err := role.RemoveMembers(ctx, client, guild, targetRole, toRemove); err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to remove members from role %s", targetRole.Name),
				err.Error(),
			)
			return nil, diags
		}
	}

	if len(toAdd) > 0 {
		if err := role.AddMembers(ctx, client, guild, targetRole, toAdd); err != nil {
			diags.AddError(
				fmt.Sprintf("Failed to add members to role %s", targetRole.Name),
				err.Error(),
			)
			return nil, diags
		}
	}

	return desired, diags
}

// unionMembersByID merges two member slices, deduplicating by user ID. Order:
// all of a (in order), then any of b not already present (in order).
func unionMembersByID(a, b []*discordgo.Member) []*discordgo.Member {
	seen := make(map[string]struct{}, len(a)+len(b))
	out := make([]*discordgo.Member, 0, len(a)+len(b))

	for _, group := range [][]*discordgo.Member{a, b} {
		for _, m := range group {
			if m == nil || m.User == nil {
				continue
			}
			if _, ok := seen[m.User.ID]; ok {
				continue
			}
			seen[m.User.ID] = struct{}{}
			out = append(out, m)
		}
	}

	return out
}

// partitionForSync computes which members need the role added and which need
// it removed, given the desired final membership and who currently holds it
// among the tracked/checked set (see SyncDeclaredMembers, FetchDeclaredHolders).
// Pure function, no API calls — this is the piece unit tests exercise most.
func partitionForSync(desired, held []*discordgo.Member) (toAdd, toRemove []*discordgo.Member) {
	desiredIDs := make(map[string]struct{}, len(desired))
	for _, m := range desired {
		if m == nil || m.User == nil {
			continue
		}
		desiredIDs[m.User.ID] = struct{}{}
	}

	heldIDs := make(map[string]struct{}, len(held))
	for _, m := range held {
		if m == nil || m.User == nil {
			continue
		}
		heldIDs[m.User.ID] = struct{}{}
		if _, wanted := desiredIDs[m.User.ID]; !wanted {
			toRemove = append(toRemove, m)
		}
	}

	for _, m := range desired {
		if m == nil || m.User == nil {
			continue
		}
		if _, has := heldIDs[m.User.ID]; !has {
			toAdd = append(toAdd, m)
		}
	}

	return toAdd, toRemove
}
