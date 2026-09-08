package discordmembers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
)

// fakeClient is a hand-written DiscordMemberClient fake. There is no existing
// mocking convention in this codebase (discordgo.Session is a concrete struct
// with no interface of its own), so this is the first.
type fakeClient struct {
	// membersByGuild backs both GuildMembersSearch (prefix-matched, mirroring
	// Discord's real "starts with" semantics) and GuildMember (exact ID match).
	membersByGuild map[string][]*discordgo.Member

	// searchErr, if set, is returned from every GuildMembersSearch call.
	searchErr error

	// getErrByKey maps "guildID/userID" to an error GuildMember should return
	// for that specific lookup instead of a normal result.
	getErrByKey map[string]error

	searchCalls int
	getCalls    int
}

func (f *fakeClient) GuildMembersSearch(guildID, query string, limit int) ([]*discordgo.Member, error) {
	f.searchCalls++
	if f.searchErr != nil {
		return nil, f.searchErr
	}
	var out []*discordgo.Member
	for _, m := range f.membersByGuild[guildID] {
		if strings.HasPrefix(m.User.Username, query) {
			out = append(out, m)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

func (f *fakeClient) GuildMember(guildID, userID string) (*discordgo.Member, error) {
	f.getCalls++
	if f.getErrByKey != nil {
		if err, ok := f.getErrByKey[guildID+"/"+userID]; ok {
			return nil, err
		}
	}
	for _, m := range f.membersByGuild[guildID] {
		if m.User.ID == userID {
			return m, nil
		}
	}
	return nil, notFoundErr()
}

// notFoundErr builds a *discordgo.RESTError shaped like a real 404 from
// GuildMember, so tests exercise the exact error-shape FetchDeclaredHolders
// switches on.
func notFoundErr() error {
	return &discordgo.RESTError{Response: &http.Response{StatusCode: http.StatusNotFound}}
}

func member(id, username string, roles ...string) *discordgo.Member {
	return &discordgo.Member{
		User:  &discordgo.User{ID: id, Username: username},
		Roles: roles,
	}
}

func TestResolveMembersByUsername_ExactMatchAmongPrefixMatches(t *testing.T) {
	client := &fakeClient{
		membersByGuild: map[string][]*discordgo.Member{
			"g1": {
				member("1", "alice"),
				member("2", "albert"),
				member("3", "alicia"),
			},
		},
	}

	got, diags := ResolveMembersByUsername(context.Background(), client, "g1", []string{"alice"})
	if diags.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", diags)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(got) != 1 || got[0].User.ID != "1" {
		t.Fatalf("expected exactly member id=1 (alice), got %+v", got)
	}
	// One Search call for the one declared username — confirms this resolves
	// per-username rather than paginating a full guild listing.
	if client.searchCalls != 1 {
		t.Fatalf("expected 1 search call, got %d", client.searchCalls)
	}
}

func TestResolveMembersByUsername_NotFoundWarnsAndSkips(t *testing.T) {
	client := &fakeClient{
		membersByGuild: map[string][]*discordgo.Member{
			"g1": {member("1", "alice")},
		},
	}

	got, diags := ResolveMembersByUsername(context.Background(), client, "g1", []string{"alice", "typo-name"})
	if diags.HasError() {
		t.Fatalf("not-found should warn, not error: %v", diags)
	}
	if len(got) != 1 || got[0].User.Username != "alice" {
		t.Fatalf("expected only alice resolved, got %+v", got)
	}
	if len(diags) != 1 {
		t.Fatalf("expected exactly one warning diagnostic, got %d: %v", len(diags), diags)
	}
	if !strings.Contains(diags[0].Summary(), "not found") {
		t.Fatalf("expected a not-found warning, got %q", diags[0].Summary())
	}
}

func TestResolveMembersByUsername_SearchErrorHardFailsButContinues(t *testing.T) {
	client := &fakeClient{searchErr: errors.New("rate limited")}

	got, diags := ResolveMembersByUsername(context.Background(), client, "g1", []string{"alice", "bob"})
	if !diags.HasError() {
		t.Fatalf("expected error diagnostics from search failure")
	}
	if len(got) != 0 {
		t.Fatalf("expected no members resolved when search fails, got %+v", got)
	}
	// Both usernames should have been attempted, not just the first.
	if client.searchCalls != 2 {
		t.Fatalf("expected 2 search attempts, got %d", client.searchCalls)
	}
}

func TestFetchDeclaredHolders_PartitionsByCurrentRole(t *testing.T) {
	client := &fakeClient{
		membersByGuild: map[string][]*discordgo.Member{
			"g1": {
				member("1", "alice", "role-a"),
				member("2", "bob"),
			},
		},
	}
	declared := []*discordgo.Member{member("1", "alice"), member("2", "bob")}

	held, diags := FetchDeclaredHolders(context.Background(), client, "g1", "role-a", declared)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if len(held) != 1 || held[0].User.ID != "1" {
		t.Fatalf("expected only alice (id=1) to hold role-a, got %+v", held)
	}
}

func TestFetchDeclaredHolders_DepartedMemberWarnsAndDrops(t *testing.T) {
	client := &fakeClient{
		membersByGuild: map[string][]*discordgo.Member{
			"g1": {member("1", "alice", "role-a")},
		},
	}
	// "bob" (id=2) is declared but no longer in the guild.
	declared := []*discordgo.Member{member("1", "alice"), member("2", "bob")}

	held, diags := FetchDeclaredHolders(context.Background(), client, "g1", "role-a", declared)
	if diags.HasError() {
		t.Fatalf("a departed member must not hard-fail Read: %v", diags)
	}
	if len(held) != 1 || held[0].User.ID != "1" {
		t.Fatalf("expected only alice to remain held, got %+v", held)
	}
	if len(diags) != 1 {
		t.Fatalf("expected exactly one warning for the departed member, got %d: %v", len(diags), diags)
	}
	if !strings.Contains(diags[0].Summary(), "left the guild") {
		t.Fatalf("expected a left-the-guild warning, got %q", diags[0].Summary())
	}
}

func TestFetchDeclaredHolders_TransientErrorHardFailsAndAborts(t *testing.T) {
	client := &fakeClient{
		membersByGuild: map[string][]*discordgo.Member{
			"g1": {member("1", "alice", "role-a"), member("2", "bob", "role-a")},
		},
		getErrByKey: map[string]error{
			// A non-404 error on the FIRST declared member must not be mistaken
			// for "member left" — it must hard-fail the whole call rather than
			// silently proceeding to check the second member and shrinking
			// membership based on a transient failure.
			"g1/1": errors.New("HTTP 429: rate limited"),
		},
	}
	declared := []*discordgo.Member{member("1", "alice"), member("2", "bob")}

	held, diags := FetchDeclaredHolders(context.Background(), client, "g1", "role-a", declared)
	if !diags.HasError() {
		t.Fatalf("expected a hard error for a non-404 failure, got %v", diags)
	}
	if held != nil {
		t.Fatalf("expected no partial result on hard failure, got %+v", held)
	}
	// Must not have proceeded to check bob after alice's transient failure.
	if client.getCalls != 1 {
		t.Fatalf("expected the loop to abort after the first hard failure, got %d GuildMember calls", client.getCalls)
	}
}

func TestPartitionForSync(t *testing.T) {
	desired := []*discordgo.Member{member("1", "alice"), member("2", "bob")}
	held := []*discordgo.Member{member("2", "bob"), member("3", "carol")}

	toAdd, toRemove := partitionForSync(desired, held)

	if len(toAdd) != 1 || toAdd[0].User.ID != "1" {
		t.Fatalf("expected alice (id=1) to need adding, got %+v", toAdd)
	}
	if len(toRemove) != 1 || toRemove[0].User.ID != "3" {
		t.Fatalf("expected carol (id=3) to need removing, got %+v", toRemove)
	}
}

func TestPartitionForSync_NoOpWhenAlreadyInSync(t *testing.T) {
	desired := []*discordgo.Member{member("1", "alice"), member("2", "bob")}
	held := []*discordgo.Member{member("1", "alice"), member("2", "bob")}

	toAdd, toRemove := partitionForSync(desired, held)

	if len(toAdd) != 0 || len(toRemove) != 0 {
		t.Fatalf("expected no changes, got toAdd=%+v toRemove=%+v", toAdd, toRemove)
	}
}

func TestUnionMembersByID_Dedupes(t *testing.T) {
	a := []*discordgo.Member{member("1", "alice"), member("2", "bob")}
	b := []*discordgo.Member{member("2", "bob"), member("3", "carol")}

	got := unionMembersByID(a, b)

	if len(got) != 3 {
		t.Fatalf("expected 3 unique members, got %d: %+v", len(got), got)
	}
	seen := map[string]bool{}
	for _, m := range got {
		if seen[m.User.ID] {
			t.Fatalf("duplicate member id=%s in union result", m.User.ID)
		}
		seen[m.User.ID] = true
	}
}

func TestSyncDeclaredMembers_UnionAndPartitionCatchRemovedMember(t *testing.T) {
	// Regression guard for the design's central correctness requirement: a
	// member present in `previous` but absent from `desired` (deleted from a
	// role's members file between applies) must actually be removed, not
	// silently kept — this is the entire reason SyncDeclaredMembers checks
	// previous ∪ desired rather than desired alone. This exercises the two
	// pure pieces SyncDeclaredMembers composes (unionMembersByID,
	// partitionForSync) together; SyncDeclaredMembers itself takes a concrete
	// *discordgo.Session (to reach go-discordutils's AddMembers/RemoveMembers)
	// and is exercised at the acceptance-test layer instead, per this
	// package's testing approach.
	const roleID = "role-a"

	previous := []*discordgo.Member{member("1", "alice"), member("2", "bob")}
	desired := []*discordgo.Member{member("1", "alice")} // bob removed from config

	tracked := unionMembersByID(previous, desired)
	if len(tracked) != 2 {
		t.Fatalf("expected the union to include bob even though he's no longer desired, got %+v", tracked)
	}

	// Simulate FetchDeclaredHolders' result: both still hold the role today.
	held := []*discordgo.Member{
		member("1", "alice", roleID),
		member("2", "bob", roleID),
	}

	toAdd, toRemove := partitionForSync(desired, held)
	if len(toAdd) != 0 {
		t.Fatalf("expected nothing to add, got %+v", toAdd)
	}
	if len(toRemove) != 1 || toRemove[0].User.ID != "2" {
		t.Fatalf("expected bob (id=2) to be removed, got %+v", toRemove)
	}
}
