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

func article(ss *dgo.Session, guild, user, article string) (content string, hidden uint64) {
	if sub, ok := state[guild].Submit[user]; ok {
		if article == "" {
			content = "Revoked \"" + sub + "\"."
			delete(state[guild].Submit, user)
		} else {
			content = "Submitted \"" + article +
					"\", revoking \"" + sub + "\"."
			state[guild].Submit[user] = article
		}
	} else if article == "" {
		content = "No article submitted."
	} else {
		content = "Submitted \"" + article + "\"."
		state[guild].Submit[user] = article
	}
	return content, 1
}

func clear(ss *dgo.Session, guild, user, null string) (content string, hidden uint64) {
	if state[guild].Host == user && state[guild].TmpHost {
		return "Temporary hosts may not clear the article list.", 1
	}

	state[guild].Submit = map[string]string{}
	return "Article list cleared.", 1
}

func host(ss *dgo.Session, guild, user, host string) (content string, hidden uint64) {
	if len(state[guild].Submit) == 0 {
		return "No articles have been submitted.", 1
	} else if state[guild].Host != "" {
		return "A round is already running.", 1
	} else if _, ok := state[guild].Submit[host]; ok {
		if host == user {
			return "You have submitted an article, which must be " +
					"revoked before you may host a round.", 1
		} else {
			return "<@" + host + "> has submitted an article, which must be " +
					"revoked before they may host a round.", 1
		}
	}

	// TODO: check whether user already has wikidt role to prevent :trolling:
	state[guild].Host, state[guild].TmpHost = host, host != user
	if host != user {
		ss.GuildMemberRoleAdd(guild, host, state[guild].Trusted)
	}

	count, ran := 0, rand.Intn(len(state[guild].Submit))
	for player, article := range state[guild].Submit {
		if count == ran {
			state[guild].Player = player
			delete(state[guild].Submit, player)
			return "A new round of wikid has begun! The article is \"" +
					article + "\", and the host is <@" +
					host + ">.", 0
		}
		count++
	}
	return
}

func guess(ss *dgo.Session, guild, user, player string) (content string, hidden uint64) {
	if state[guild].Host == "" {
		return "A round is not currently running.", 1
	} else if user == state[guild].Host && player != user {
		if player == state[guild].Player {
			content = "<@" + player + "> was guessed to have submitted " +
					"the article, which is correct."
		} else {
			content = "<@" + player + "> was guessed to have submitted the article, " +
					"but it was actually <@" + state[guild].Player + ">."
		}
	} else if user != player {
		return "You must ping yourself to end the round prematurely.", 1
	} else {
		content = "The round has been ended prematurely."
	}

	if state[guild].TmpHost {
		ss.GuildMemberRoleRemove(guild, state[guild].Host, state[guild].Trusted)
	}
	state[guild].Host, state[guild].Player = "", ""
	return content, 0
}

func ban(ss *dgo.Session, guild, user, player string) (content string, hidden uint64) {
	if user == state[guild].Host && state[guild].TmpHost {
		return "Temporary hosts may not ban users.", 1
	}

	delete(state[guild].Submit, player)
	if ss.GuildMemberRoleAdd(guild, player, state[guild].Banned) == nil {
		return "<@" + player + "> has been banned.", 0
	}
	return "Unable to ban <@" + player + ">.", 1
}
