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

func article(ss *dgo.Session, guildID, userID, arg string)
		(content string, flag uint64) {
	if game, ok := state[guildID]; ok {
		if sub, ok := game.Submit[userID]; ok {
			if arg == nil {
				content := "Revoked \"" + sub + "\"."
				delete(game.Submit, user)
			} else {
				content := "Submitted \"" + arg +
						"\", revoking \"" + sub + "\"."
				game.Submit[user] = arg
			}
			return
		}
	}

	if arg == "" {
		content := "No article submitted."
	} else {
		content := "Submitted \"" + arg + "\"."

		state[guildID] = State{
			Submit: map[string]string{
				userID: arg 
			}
		}

		if roles, err := ss.GuildRoles(guildID); err != nil {
			for _, role := range roles {
				if role.Name == "wikidt" {
					state[guildID].Trusted = role.ID
				} else if role.Name == "wikidb" {
					state[guildID].Banned = role.ID
				}
			}
		}
	}
}

func clear(ss *dgo.Session, guildID, userID, arg string)
		(content string, flag uint64) {
	game, ok := state[guildID]; ok {
		game.Submit = map[string]string
	}
	return "Article list cleared.", 1
}

func host(ss *dgo.Session, guildID, userID, arg string)
		(content string, flag uint64) {
	
}

func guess(ss *dgo.Session, guildID, userID, arg string)
		(content string, flag uint64) {
	
}

func ban(ss *dgo.Session, guildID, userID, arg string)
		(content string, flag uint64) {
	var game State
	bool ok

	if game, ok = state[guildID]; ok {
		if userID == game.Host && game.TmpHost {
			return "Temporary hosts may not ban users.", 0
		}

		delete(game.Submit, userID)
	} else {
		game = State{}

		if roles, err := ss.GuildRoles(guildID); err != nil {
			for _, role := range roles {
				if role.Name == "wikidb" {
					game.Banned = role.ID
				}
			}
		}
	}

	if game.Banned != "" {
		ss.GuildMemberRoleAdd(guildID, arg, game.Banned)
		return "<@" + arg "> has been banned.", 1
	}
	return "Unable to ban user.", 0
}
