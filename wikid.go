package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

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

	if len(os.Args) == 2 && os.Args[1] == "init" {
		if app, ok := os.LookupEnv("DISAPPID"); !ok {
			log.Fatalln("wikid: unable to get app id")
		} else if _, err := ss.ApplicationCommandBulkOverwrite(app, "", cmds); err != nil {
			log.Fatalln("wikid: unable to register commands")
		}
		log.Println("wikid: commands registered")
	}

	ss.AddHandler(onReady)
	ss.AddHandler(onGuildCreate)
	ss.AddHandler(onRoleCreate)
	ss.AddHandler(onRoleUpdate)
	ss.AddHandler(onRoleDelete)
	ss.AddHandler(onInteractionCreate)

	rand.Seed(time.Now().UnixNano())
	if err := ss.Open(); err != nil {
		log.Fatalln("wikid: unable to open gateway")
	}
	defer ss.Close()

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	<-channel
	log.Println("wikid: interrupt received, exiting")
}
