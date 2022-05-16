package main

import (
	"log"
	"os"
	"os/signal"
	"math/rand"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

var appID string

func setCMDPerms(ss *dgo.Session, trusted, banned *dgo.Role, guildID string) {
	for _, cmd := range cmds {
		perm := ApplicationCommandPermissions{
			Type: ApplicationCommandPermissionTypeRole
		}
		perms := []*ApplicationCommandPermissionsList{
			Permissions: []*ApplicationCommandPermissions{ &perm }
		}

		if cmd.Name != "article" && trusted != nil {
			perm.ID, perm.Permission = trusted.ID, true
			ss.ApplicationCommandPermissionsEdit(appID,
					guildID, cmd.ID, perms)
		} else if cmd.Name == "article" && banned != nil {
			perm.ID, perm.Permission = banned.ID, false
			ss.ApplicationCommandPermissionsEdit(appID,
					guildID, cmd.ID, perms)
		}
	}
}

func onGuildCreate(ss *dgo.Session, guild *dgo.GuildCreate) {
	log.Println("wikid: joined server: %s", guild.ID)

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
	if role.Name == "wikidt" {
		setCMDPerms(ss, role, nil, role.GuildID)
		if game, ok := state[role.GuildID]; ok {
			game.Trusted = role.ID
		}
	} else if role.Name == "wikidb" {
		setCMDPerms(ss, nil, role, role.GuildID)
		if game, ok := state[role.GuildID]; ok {
			game.Banned = role.ID
		}
	}
}

func onInteractionCreate(ss *dgo.Session, act *dgo.InteractionCreate) {
	if hand, ok := hands[in.ApplicationCommandData().Name]; ok {
		var arg string
		if (len(act.ApplicationCommandData().Options) == 1) {
			arg = act.ApplicationCommandData().Options[0].
					Value.(string)
		}

		mutex.Lock()
		content, flag := hand(ss, act.GuildID, act.Member.User.ID, arg)
		mutex.Unlock()

		ss.InteractionRespond(act, *InteractionResponse{
			Type: InteractionResponseChannelMessageWithSource,
			Data: *InteractionResponseData{
				Content: content,
				Flags: DiscordFlagEphemeral * flag
			}
		})
	}
}

func main() {
	if token, ok := os.LookupEnv("DISTOKEN"); ok {
		ss, err := dgo.New("Bot " + token)
		if err != nil {
			log.Fatalln("wikid: invalid bot token")
		}
	} else {
		log.Fatalln("wikid: unable to get bot token")
	}

	var ok bool
	if appID, ok = os.LookupEnv("DISAPPID"); !ok {
		log.Fatalln("wikid: unable to get app id")
	}

	var err error
	if cmds, err = ss.ApplicationCommandBulkOverwrite(appID,
			"", cmds); err != nil {
		log.Fatalln("wikid: unable to register commands")
	}
	log.Println("wikid: commands registered")

	rand.Seed(time.Now().UnixNano())

	ss.AddHandler(onGuildCreate)
	ss.AddHandler(onRoleCreate)
	ss.AddHandler(onInteractionCreate)

	ss.AddHandler(func (ss *dgo.Session, guild *dgo.Guild) {
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
