package discordmembers

import "github.com/bwmarrin/discordgo"

// DiscordMemberClient is the narrow surface this package needs from a Discord
// REST client for member resolution and role-membership checks. It exists so the
// resolution/reconciliation logic in discord.go can be unit-tested against a
// hand-written fake, without requiring discordgo.Session (a concrete struct with
// no interface of its own, and no existing mocking precedent anywhere in this
// codebase) to be mocked.
//
// It deliberately covers only ungated Discord REST endpoints:
//   - GuildMembersSearch: GET /guilds/{id}/members/search
//   - GuildMember: GET /guilds/{id}/members/{user} (single member — not the
//     gated List Guild Members endpoint)
//
// Nothing in this package should call Session.GuildMembers (List Guild Members).
// That endpoint requires the GUILD_MEMBERS privileged intent even over REST —
// this whole package exists to stop depending on it. See discord.go.
type DiscordMemberClient interface {
	GuildMembersSearch(guildID, query string, limit int) ([]*discordgo.Member, error)
	GuildMember(guildID, userID string) (*discordgo.Member, error)
}

// sessionMemberClient adapts a *discordgo.Session to DiscordMemberClient,
// dropping the variadic RequestOption parameters that a plain interface method
// set can't express. It touches nothing else about the session (state cache,
// rate limiter, etc.) — this is a signature-narrowing wrapper only, constructed
// inline at each call site that needs it.
type sessionMemberClient struct {
	session *discordgo.Session
}

// NewSessionMemberClient wraps a live discordgo.Session for use wherever this
// package needs a DiscordMemberClient.
func NewSessionMemberClient(session *discordgo.Session) DiscordMemberClient {
	return &sessionMemberClient{session: session}
}

func (c *sessionMemberClient) GuildMembersSearch(guildID, query string, limit int) ([]*discordgo.Member, error) {
	return c.session.GuildMembersSearch(guildID, query, limit)
}

func (c *sessionMemberClient) GuildMember(guildID, userID string) (*discordgo.Member, error) {
	return c.session.GuildMember(guildID, userID)
}
