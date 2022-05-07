package main

import (
	"log"
	"os"
	"os/signal"
	"math/rand"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

type NRoles struct {
	Trusted *dgo.Role
	Banned  *dgo.Role
}

var appID string

func setCMDPerms(ss *dgo.Session, roles *NRoles, guildID string) {
	for _, cmd := range cmds {
		perm := ApplicationCommandPermissions{
			Type: ApplicationCommandPermissionTypeRole
		}
		perms := []*ApplicationCommandPermissionsList{
			Permissions: []*ApplicationCommandPermissions{ &perm }
		}

		if cmd.Name != "article" && roles.Trusted != nil {
			perm.ID = roles.Trusted.ID; perm.Permission = true
			ss.ApplicationCommandPermissionsEdit(appID,
					guildID, cmd.ID, perms)
		} else if cmd.Name == "article" && roles.Banned != nil {
			perm.ID = roles.Banned.ID; perm.Permission = false
			ss.ApplicationCommandPermissionsEdit(appID,
					guildID, cmd.ID, perms)
		}
	}
}

func onReady(ss *dgo.Session, guild *dgo.Guild) {
	log.Println("wikid: opened gateway")

	var err error
	if cmds, err = ss.ApplicationCommandBulkOverwrite(appID, "", cmds);
			err != nil {
		log.Fatalln("wikid: unable to register commands")
	}

	log.Println("wikid: commands registered")
}

func onGuildCreate(ss *dgo.Session, guild *dgo.GuildCreate) {
	log.Println("wikid: joined server: %s", guild.ID)

	var roles NRoles
	for _, role := range guild.NRoles {
		if roles.Trusted != nil && role.Name == "wikidt" {
			roles.Trusted = role
		} else if roles.Banned != nil && role.Name == "wikidb" {
			roles.Banned = role
		}
	}
	if roles.Trusted == nil && roles.Banned == nil {
		return
	}
	setCMDPerms(ss, &roles, guild.ID)
}

func onRoleCreate(ss *dgo.Session, role *dgo.GuildRoleCreate) {
	var roles NRoles
	if role.Name == "wikidt" {
		roles.Trusted = role
	} else if role.Name == "wikidb" {
		roles.Banned = role
	} else {
		return
	}
	setCMDPerms(ss, &roles, role.GuildID)
}

func onInteractionCreate(ss *dgo.Session, act *dgo.InteractionCreate) {
	if hand, ok := hands[in.ApplicationCommandData().Name]; ok {
		mutex.Lock()
		hand(ss. act)
		mutex.Unlock()
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

	rand.Seed(time.Now().UnixNano())

	ss.AddHandler(onReady)
	ss.AddHandler(onGuildCreate)
	ss.AddHandler(onRoleCreate)
	ss.AddHandler(onInteractionCreate)

	if err := ss.Open(); err != nil {
		log.Fatalln("wikid: unable to open gateway")
	}
	defer ss.Close()

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	<-channel
	log.Println("wikid: interrupt received, exiting")
}
