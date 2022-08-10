package main

import (
	"flag"
	"github.com/peam1146/mcv-notifier/src/controller"
	"github.com/peam1146/mcv-notifier/src/database"
	"github.com/peam1146/mcv-notifier/src/service"
	"log"
)

var (
	ChannelID    = flag.String("channel", "", "Channel IDx")
	BotToken     = flag.String("token", "", "Bot access token")
	Username     = flag.String("username", "", "Username")
	Password     = flag.String("password", "", "Password")
	DatabasePath = flag.String("database", "database.db", "Database path")
	Debug        = flag.Bool("debug", false, "Debug mode")
)

func init() { flag.Parse() }

func main() {

	db := database.NewSqliteDatabase(*DatabasePath, *Debug)
	ns := service.NewNotificationService(db)
	scrapeService := service.NewScraperService(ns, *Username, *Password)

	discordController := controller.NewDiscordController(*BotToken, *ChannelID, scrapeService)
	defer discordController.Close()

	wait := make(chan bool, 0)
	discordController.Init(wait)

	if err := discordController.Open(); err != nil {
		log.Fatal(err)
	}

	<-wait
	log.Println("Start scraping...")
	discordController.SendNotification()
	log.Println("Scraping finished")
}
