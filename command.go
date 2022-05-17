package main

import dgo "github.com/bwmarrin/discordgo"

func makeBool(data bool) *bool {
	return &data
}

var cmds = []*dgo.ApplicationCommand{
	{
		Name: "article",
		Description: "Submit an article to possibly be selected next round.",
		Options: []*dgo.ApplicationCommandOption{
			{
				Type: dgo.ApplicationCommandOptionString,
				Name: "article",
				Description: "The article to submit. Supplying no article will " +
						"revoke a previous submission.",
			},
		},
	},
	{
		Name: "clear",
		Description: "Clear the list of articles.",
		DefaultPermission: makeBool(false),
	},
	{
		Name: "host",
		Description: "Host a round of wikid by randomly selecting an article.",
		Options: []*dgo.ApplicationCommandOption{
			{
				Type: dgo.ApplicationCommandOptionUser,
				Name: "host",
				Description: "The user to host the round. Supplying no user " +
						"will default to self.",
			},
		},
		DefaultPermission: makeBool(false),
	},
	{
		Name: "guess",
		Description: "End a round of wikid by guessing who submitted the randomly " +
				"selected article.",
		Options: []*dgo.ApplicationCommandOption{
			{
				Type: dgo.ApplicationCommandOptionUser,
				Name: "player",
				Description: "The player who submitted the article. Enter " +
						"yourself to end the round early.",
				Required: true,
			},
		},
		DefaultPermission: makeBool(false),
	},
	{
		Name: "ban",
		Description: "Ban a user from submitting articles.",
		Options: []*dgo.ApplicationCommandOption{
			{
				Type: dgo.ApplicationCommandOptionUser,
				Name: "player",
				Description: "The player to ban from submitting articles",
				Required: true,
			},
		},
		DefaultPermission: makeBool(false),
	},
}
