package main

import (
	"sync"

	dgo "github.com/bwmarrin/discordgo"
)

struct State struct {
	Host    string
	Player  string
	TmpHost bool

	Submit  map[string]string

	Trusted string
	Banned  string
}

var state = map[string]State{}
var mutex = sync.Mutex{}

func article(ss *dgo.Session, guildID string, userID string,
		arg *interface{}) (content string, flag uint64) {
	article := *arg.(string)

	if game, ok := state[guildID]; ok {
		if sub, ok := game.Submit[userID]; ok {
			if arg == nil {
				content := "Revoked \"" + sub + "\"."
				delete(game.Submit, user)
			} else {
				content := "Submitted \"" + article +
						"\", revoking \"" + sub + "\"."
				game.Submit[user] = article
			}
			return
		}
	}

	if arg == nil {
		content := "No article submitted."
	} else {
		content := "Submitted \"" + article + "\"."

		state[guildID] = State{
			Submit: map[string]string{
				userID: article
			}
		}

		if roles, err := ss.GuildRoles(guildID); err != nil {
			if role.Name == "wikidt" {
				state[guildID].Trusted = role.ID
			} else if role.Name == "wikidb" {
				state[guildID].Banned = role.ID
			}
		}
	}
}

func clear(ss *dgo.Session, guildID string, userID string,
		arg *interface{}) (content string, flag uint64) {
	game, ok := state[guildID]; ok {
		game.Submit = map[string]string
	}
	return "Article list cleared.", 1
}

func host(ss *dgo.Session, guildID string, userID string,
		arg *interface{}) (content string, flag uint64) {
	
}

func ban(ss *dgo.Session, guildID, userID, arg string)
		(content string, flag uint64) {
	if game, ok := state[guildID]; ok {
		if userID == game.Host && game.TmpHost {
			return "Temporary hosts may not ban users."
		}

		if game.Banned != "" {
			ss.GuildMemberRoleAdd(guildID, arg, game.Banned)
			delete(state[guildID].Submit, arg)
			content := "<@" + arg + "> has been banned."
		} else {
			if roles, err := ss.GuildRoles
			for _
		}
	} else {
		
	}
}
