package main

import dgo "github.com/bwmarrin/discordgo"

var cmds = []*ApplicationCommand{
	{
		Name: "article",
		Description: "Submit an article to possibly be selected " +
				"next round.",
		Options: []*ApplicationCommandOption{
			Type: ApplicationCommandOptionString,
			Name: "article",
			Description: "The article to submit. Supplying no " +
					"article will revoke a previous " +
					"submission."
		},
	},
	{
		Name: "clear",
		Description: "Clear the list of articles.",
		DefaultPermissions: false
	},
	{
		Name: "host",
		Description: "Host a round of wikid by randomly selecting " +
				"an article.",
		Options: []*ApplicationCommandOption{
			Type: ApplicationCommandOptionUser,
			Name: "host",
			Description: "The user to host the round. Supplying " +
					"no user will default to self."
		},
		DefaultPermissions: false
	},
	{
		Name: "guess",
		Description: "End a round of wikid by guessing who " +
				"submitted the randomly selected article.",
		Options: []*ApplicationCommandOption{
			Type: ApplicationCommandOptionUser,
			Name: "player",
			Description: "The player who submitted the article. " +
					"Enter yourself to end the round " +
					"early.",
			Required: true
		},
		DefaultPermissions: false
	},
	{
		Name: "ban",
		Description: "Ban a user from submitting articles.",
		Options: []*ApplicationCommandOption{
			Type: ApplicationCommandOptionUser,
			Name: "player",
			Description: "The player to ban from submitting " +
					"articles",
			Required: true
		},
		DefaultPermissions: false
	}
}
