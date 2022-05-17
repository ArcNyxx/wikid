package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

var app string
var mutex = sync.Mutex{}

func setCMDPerms(ss *dgo.Session, trusted, banned *dgo.Role, guild string) {
	for _, cmd := range cmds {
		perms := &dgo.ApplicationCommandPermissionsList{
			Permissions: []*dgo.ApplicationCommandPermissions{
				{
					Type: dgo.ApplicationCommandPermissionTypeRole,
				},
			},
		}

		if cmd.Name != "article" && trusted != nil {
			perms.Permissions[0].ID, perms.Permissions[0].Permission = trusted.ID, true
		} else if cmd.Name == "article" && banned != nil {
			perms.Permissions[0].ID, perms.Permissions[0].Permission = banned.ID, false
		}

		ss.ApplicationCommandPermissionsEdit(app, guild, cmd.ID, perms)
	}
}

func onGuildCreate(ss *dgo.Session, guild *dgo.GuildCreate) {
	log.Println("wikid: joined server:", guild.ID)

	var trusted, banned *dgo.Role
	for _, role := range guild.Roles {
		if role.Name == "wikidt" {
			trusted = role
		} else if role.Name == "wikidb" {
			banned = role
		}
	}
	if trusted != nil || banned != nil {
		setCMDPerms(ss, trusted, banned, guild.ID)
	}
}

func onRoleCreate(ss *dgo.Session, role *dgo.GuildRoleCreate) {
	if role.Role.Name == "wikidt" {
		setCMDPerms(ss, role.Role, nil, role.GuildID)
		if game, ok := state[role.GuildID]; ok {
			game.Trusted = role.Role.ID
		}
	} else if role.Role.Name == "wikidb" {
		setCMDPerms(ss, nil, role.Role, role.GuildID)
		if game, ok := state[role.GuildID]; ok {
			game.Banned = role.Role.ID
		}
	}
}

func onInteractionCreate(ss *dgo.Session, act *dgo.InteractionCreate) {
	hands := map[string]func(ss *dgo.Session, guild, user, article string) (content string,
			flag bool){
		"article": article, "clear": clear, "host": host, "guess": guess, "ban": ban,
	}

	if hand, ok := hands[act.ApplicationCommandData().Name]; ok {
		var arg string
		if (len(act.ApplicationCommandData().Options) == 1) {
			arg = act.ApplicationCommandData().Options[0].
					Value.(string)
		}

		mutex.Lock()
		content, flag := hand(ss, act.GuildID, act.Member.User.ID, arg)
		mutex.Unlock()

		var flags uint64
		if !flag {
			flags = uint64(dgo.MessageFlagsEphemeral)
		}

		ss.InteractionRespond(act.Interaction, &dgo.InteractionResponse{
			Type: dgo.InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: content,
				Flags: flags,
			},
		})
	}
}

func main() {
	var ss *dgo.Session

	if token, ok := os.LookupEnv("DISTOKEN"); ok {
		var err error
		if ss, err = dgo.New("Bot " + token); err != nil {
			log.Fatalln("wikid: invalid bot token")
		}
	} else {
		log.Fatalln("wikid: unable to get bot token")
	}

	var ok bool
	if app, ok = os.LookupEnv("DISAPPID"); !ok {
		log.Fatalln("wikid: unable to get app id")
	}

	var err error
	if cmds, err = ss.ApplicationCommandBulkOverwrite(app,
			"", cmds); err != nil {
		log.Fatalln("wikid: unable to register commands")
	}
	log.Println("wikid: commands registered")

	rand.Seed(time.Now().UnixNano())

	ss.AddHandler(onGuildCreate)
	ss.AddHandler(onRoleCreate)
	ss.AddHandler(onInteractionCreate)

	ss.AddHandler(func(ss *dgo.Session, ready *dgo.Ready) {
		log.Println("wikid: opened gateway")
	})

	if err := ss.Open(); err != nil {
		log.Fatalln("wikid: unable to open gateway")
	}
	defer ss.Close()

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	<-channel
	log.Println("wikid: interrupt received, exiting")
}
