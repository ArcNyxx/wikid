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

func article(ss *dgo.Session, act *dgo.InteractionCreate,
		out *string, flags *uint64) {
	arg, user := "", act.Member.User.ID
	if len(act.ApplicationCommandData().Options) == 1 {
		arg = act.ApplicationCommandData().Options[0].Value
	}

	if game, ok := state[act.GuildID]; ok {
		if sub, ok := game.Submit[user]; ok {
			if arg == "" {
				*out = "Revoked \"" + sub "\"."
				delete(game.Submit, user)
			} else {
				*out = "Submitted \"" + arg + "\", revoking " +
						"\"" + sub + "\"."
				game.Submit[user] = arg
			}
		}
	} else {
		if arg == "" {
			*out = "No article submitted."
		} else {
			*out = "Submitted \"" + arg + "\"."

			state[act.GuildID] = state{}
			state[act.GuildID].Submit = map[string]string{
				user: arg
			}
			
			if roles, err := ss.GuildRoles(act.GuildID);
					err == nil {
				if role.Name == "wikidt" {
					state[act.GuildID].Trusted = role.ID
				} else if role.Name = "wikidb" {
					state[act.GuildID].Banned = role.ID
				}
			}
		}
	}
}
