package member

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// FetchById fetches a member by ID.
func FetchById(ctx context.Context, client *discordgo.Session, guild_id, member_id string) (*discordgo.Member, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		member, err := client.GuildMember(guild_id, member_id)
		if err != nil {
			return nil, err
		}

		return member, nil
	}
}

// FetchByName fetches a member by name.
func FetchByName(ctx context.Context, client *discordgo.Session, guild_id, name string) (*discordgo.Member, error) {
	members, err := Fetch(ctx, client, guild_id, 1000)
	if err != nil {
		return nil, err
	}

	for _, m := range members {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if m.User.Username == name {
				return m, nil
			}
		}
	}

	// If the member is not found in the first 1000 members, we need to fetch the rest of the members
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			last_member := members[len(members)-1]
			members, err = client.GuildMembers(guild_id, last_member.User.ID, 1000)
			if err != nil {
				return nil, err
			}

			for _, m := range members {
				if m.User.Username == name {
					return m, nil
				}
			}

			if len(members) < 1000 {
				// We have reached the end of the members list as the last fetch returned less than 1000 members
				return nil, fmt.Errorf("member not found: name=%s", name)
			}
		}
	}
}

// Fetch fetches members in a guild.
func Fetch(ctx context.Context, client *discordgo.Session, guild_id string, count int) ([]*discordgo.Member, error) {
	if count <= 0 {
		count = 1000 // Maximum number of members that can be fetched at once
	}

	var members []*discordgo.Member

	lastMember := ""
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			fetched, err := client.GuildMembers(guild_id, lastMember, count)
			if err != nil {
				return nil, err
			}

			members = append(members, fetched...)

			if len(fetched) < count {
				// No more members to fetch
				return members, nil
			}

			lastMember = fetched[len(fetched)-1].User.ID
		}
	}
}

// All fetches all members in a guild.
func All(ctx context.Context, client *discordgo.Session, guild_id string) ([]*discordgo.Member, error) {
	members, err := Fetch(ctx, client, guild_id, -1)
	if err != nil {
		return nil, err
	}

	return members, nil
}
