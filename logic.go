package main

import (
	"math/rand"

	dgo "github.com/bwmarrin/discordgo"
)

type State struct {
	Host    string
	Player  string
	TmpHost bool

	Submit  map[string]string

	Trusted string
	Banned  string
}

var state = map[string]*State{}

func article(ss *dgo.Session, guild, user, article string) (content string, flag bool) {
	if game, ok := state[guild]; ok {
		if prev, ok := game.Submit[user]; ok {
			if article == "" {
				content = "Revoked \"" + prev + "\"."
				delete(game.Submit, user)
			} else {
				content = "Submitted \"" + article +
						"\", revoking \"" + prev + "\"."
				game.Submit[user] = article
			}
		} else if article == "" {
			content = "No article submitted."
		} else {
			content = "Submitted \"" + article + "\"."
			game.Submit[user] = article
		}
	} else if article == "" {
		content = "No article submitted."
	} else {
		content = "Submitted \"" + article + "\"."

		state[guild] = &State{
			Submit: map[string]string{
				user: article,
			},
		}

		if roles, err := ss.GuildRoles(guild); err == nil {
			for _, role := range roles {
				if role.Name == "wikidt" {
					state[guild].Trusted = role.ID
				} else if role.Name == "wikidb" {
					state[guild].Banned = role.ID
				}
			}
		}
	}
	return
}

func clear(ss *dgo.Session, guild, user, null string) (content string, flag bool) {
	if game, ok := state[guild]; ok {
		if game.Host == user && game.TmpHost {
			return "Temporary hosts may not clear the article list.", false
		}
		game.Submit = map[string]string{}
	}
	return "Article list cleared.", true
}

func host(ss *dgo.Session, guild, user, host string) (content string, flag bool) {
	if host == "" {
		host = user
	}

	if game, ok := state[guild]; !ok || len(game.Submit) == 0 {
		content = "No articles have been submitted."
	} else if game.Host != "" {
		content = "A round is already running."
	} else if _, ok := game.Submit[host]; ok {
		if host == user {
			content = "You have submitted an article, which must be revoked before " +
					"you may host a round."
		} else {
			content = "<@" + host + "> has submitted an article, which must be " +
					"revoked before they may host a round."
		}
	} else {
		game.Host = host
		game.TmpHost = host != user
		if host != user {
			ss.GuildMemberRoleAdd(guild, host, game.Trusted)
		}

		count, ran := 0, rand.Intn(len(game.Submit))
		for player, article := range game.Submit {
			if count == ran {
				game.Player = player
				return "A new round of wikid has begun! The article is \"" + 
						article + "\" and the host is <@" + host + ">.", true
			}
		}
	}
	return
}

func guess(ss *dgo.Session, guild, user, player string) (content string, flag bool) {
	game, ok := state[guild]

	if !ok || game.Host == "" {
		return "A round is not currently running.", false
	} else if user == game.Host {
		if player == game.Player {
			content = "<@" + player + "> was guessed to have submitted the article, " +
					"which is correct."
		} else {
			content = "<@" + player + "> was guessed to have submitted the article, " +
					"but it was actually <@" + game.Player + ">."
		}
	} else if user != player {
		return "You must ping yourself to end the round prematurely.", false
	} else {
		content = "The round has been ended prematurely."
	}

	if game.TmpHost {
		ss.GuildMemberRoleRemove(guild, game.Host, game.Trusted)
	}
	game.Host, game.Player = "", ""
	return
}

func ban(ss *dgo.Session, guild, user, player string) (content string, flag bool) {
	game, ok := state[guild]
	if ok {
		if user == game.Host && game.TmpHost {
			return "Temporary hosts may not ban users.", false
		}
		delete(game.Submit, user)
	} else {
		game = &State{}

		if roles, err := ss.GuildRoles(guild); err == nil {
			for _, role := range roles {
				if role.Name == "wikidb" {
					game.Banned = role.ID
				}
			}
		}
	}

	if game.Banned != "" {
		if ss.GuildMemberRoleAdd(guild, player, game.Banned) == nil {
			return "<@" + player + "> has been banned.", true
		}
	}

	return "Unable to ban user.", false
}
