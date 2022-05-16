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

var state = map[string]*State{}
var mutex = sync.Mutex{}

func article(ss *dgo.Session, guildID, userID, arg string)
		(content string, flag uint64) {
	if game, ok := state[guildID]; ok {
		if sub, ok := game.Submit[userID]; ok {
			if arg == "" {
				content = "Revoked \"" + sub + "\"."
				delete(game.Submit, user)
			} else {
				content = "Submited \"" + arg + "\", " +
						"revoking \"" + sub + "\"."
				game.Submit[user] = arg
			}
		} else if arg == "" {
			content = "No article submitted."
		} else {
			content = "Submitted \"" + arg "\"."
			game.Submit[user] = arg
		}
	} else if arg == "" {
		content = "No article submitted."
	} else {
		content = "Submitted \"" + arg + "\"."

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
	if game, ok := state[guildID]; ok {
		if game.Host == userID && game.TmpHost {
			return "Temporary hosts may not clear " +
					"the article list.", 0
		}
		game.Submit = map[string]string{}
	}
	return "Article list cleared.", 1
}

func host(ss *dgo.Session, guildID, userID, arg string)
		(content string, flag uint64) {
	if game, ok := state[guildID]; !ok || len(game.Submit) == 0 {
		content = "No articles have been submitted."
	} else if game.Host != 0 {
		content = "A round is already running."
	} else if _, ok := game.Submit[arg]; ok {
		if arg == userID {
			content = "You may not host a round and " +
					"have submitted an article."
		} else {
			content = "<@" + arg + "> has submitted an article, " +
					"which must be revoked before they " +
					"may host a round."
		}
	} else {
		game.TmpHost = userID == arg
		if (userID != arg) {
			ss.GuildMemberRoleAdd(guildID, arg, game.Trusted)
		}

		count, ran := 0, rand.Intn(len(game.Submit))
		for player, article := range game.Submit {
			if count == ran {
				game.Player = player
				if arg == "" {
					game.Host = userID
				} else {
					game.Host = arg
				}

				return "A new round of wikid has begun! " +
						"The article is \"" + article +
						"\", and the host is " +
						"<@" + player + ">.", 1
			}
		}
	}
}

func guess(ss *dgo.Session, guildID, userID, arg string)
		(content string, flag uint64) {
	if game, ok := state[guildID]; !ok || game.Host == "" {
		return "A round is not currently running.", 0
	} else if userID == game.Host {
		if arg == game.Player {
			content = "<@" + arg + "> was guessed to have " +
					"submitted the article, and is correct."
		} else {
			content = "<@" + arg + "> was guessed to have " +
					"submitted the article, but it was " +
					"actually <@" + game.Player + ">."
		}
	} else if userID != arg {
		return "You must ping yourself to end the " +
				"round prematurely.", 0
	} else {
		content = "The round has been ended prematurely."
	}
	flag = 1

	if game.TmpHost {
		ss.GuildMemberRoleRemove(guildID, game.Host, game.Trusted)
	}
	game.Host, game.Player = "", ""
}

func ban(ss *dgo.Session, guildID, userID, arg string)
		(content string, flag uint64) {
	game, ok := state[guildID]
	if ok && userID == game.Host && game.TmpHost {
		return "Temporary hosts may not ban users.", 0
	} else if ok {
		delete(game.Submit, userID)
	} else {
		game = *State{}

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
		return "<@" + arg + "> has been banned.", 1
	}
	return "Unable to ban user.", 0
}
