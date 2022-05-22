package main

import (
	"log"
	"sync"

	dgo "github.com/bwmarrin/discordgo"
)

var mutex = sync.Mutex{}

func onReady(ss *dgo.Session, ready *dgo.Ready) {
	log.Println("wikid: opened gateway")
}

func onGuildCreate(ss *dgo.Session, guild *dgo.GuildCreate) {
	log.Println("wikid: joined server:", guild.ID)

	mutex.Lock()
	state[guild.ID] = &State{
		Submit: map[string]string{},
	}
	if roles, err := ss.GuildRoles(guild.ID); err == nil {
		for _, role := range roles {
			if role.Name == "wikidt" {
				state[guild.ID].Trusted = role.ID
			} else if role.Name == "wikidb" {
				state[guild.ID].Banned = role.ID
			}
		}
	}
	mutex.Unlock()
}

func onRoleCreate(ss *dgo.Session, role *dgo.GuildRoleCreate) {
	mutex.Lock()
	if role.Role.Name == "wikidt" {
		state[role.GuildID].Trusted = role.Role.ID
	} else if role.Role.Name == "wikidb" {
		state[role.GuildID].Banned = role.Role.ID
	}
	mutex.Unlock()
}

func onRoleUpdate(ss *dgo.Session, role *dgo.GuildRoleUpdate) {
	mutex.Lock()
	if role.Role.ID == state[role.GuildID].Trusted {
		state[role.GuildID].Trusted = ""
	} else if role.Role.ID == state[role.GuildID].Banned {
		state[role.GuildID].Banned = ""
	}

	if role.Role.Name == "wikidt" {
		state[role.GuildID].Trusted = role.Role.ID
	} else if role.Role.Name == "wikidb" {
		state[role.GuildID].Banned = role.Role.ID
	}
	mutex.Unlock()
}

func onRoleDelete(ss *dgo.Session, role *dgo.GuildRoleDelete) {
	mutex.Lock()
	if role.RoleID == state[role.GuildID].Trusted {
		state[role.GuildID].Trusted = ""
	} else if role.RoleID == state[role.GuildID].Banned {
		state[role.GuildID].Banned = ""
	}
	mutex.Unlock()
}

func findRole(find string, roles []string) bool {
	for _, role := range roles {
		if find == role {
			return true
		}
	}
	return false
}

func onInteractionCreate(ss *dgo.Session, act *dgo.InteractionCreate) {
	content, flags := "You do not have adequate permissions.", uint64(1)
	handles := map[string]func(*dgo.Session, string, string, string) (string, uint64){
		"clear": clear, "host": host, "guess": guess, "ban": ban,
	}

	mutex.Lock()
	if member, err := ss.GuildMember(act.GuildID, act.Member.User.ID); err == nil {
		var arg string
		if len(act.ApplicationCommandData().Options) == 1 {
			arg = act.ApplicationCommandData().Options[0].Value.(string)
		}

		if act.ApplicationCommandData().Name == "article" {
			if state[act.GuildID].Banned == "" ||
					!findRole(state[act.GuildID].Banned, member.Roles) {
				content, flags = article(ss, act.GuildID, act.Member.User.ID, arg)
			}
		} else {
			if state[act.GuildID].Trusted != "" &&
					findRole(state[act.GuildID].Trusted, member.Roles) {
				content, flags = handles[act.ApplicationCommandData().Name](ss,
						act.GuildID, act.Member.User.ID, arg)
			}
		}
	}
	mutex.Unlock()

	flags *= uint64(dgo.MessageFlagsEphemeral)
	ss.InteractionRespond(act.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: content,
			Flags: flags,
		},
	})
}
