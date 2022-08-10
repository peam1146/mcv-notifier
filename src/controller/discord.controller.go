package controller

import (
	"github.com/bwmarrin/discordgo"
	"github.com/peam1146/mcv-notifier/src/service"
	"log"
)

type DiscordController interface {
	SendNotification()
	Init() *discordgo.Session
	GetSession() *discordgo.Session
	Open() error
	Close()
}

type discordController struct {
	session       *discordgo.Session
	channelID     string
	scrapeService service.ScraperService
}

func NewDiscordController(token, channelID string, ss service.ScraperService) *discordController {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	return &discordController{discord, channelID, ss}
}

func (ds *discordController) Init(wait chan bool) *discordgo.Session {
	discord := ds.session

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		wait <- true
	})
	return discord
}

func (ds *discordController) GetSession() *discordgo.Session {
	return ds.session
}

func (ds *discordController) Open() error {
	return ds.session.Open()
}

func (ds *discordController) Close() {
	ds.session.Close()
}

func (ds *discordController) SendNotification() {
	notifications, err := ds.scrapeService.Scraping()
	if err != nil {
		ds.session.ChannelMessageSend(ds.channelID, "Error: "+err.Error())
	}
	for _, notification := range notifications {
		ds.session.ChannelMessageSendEmbed(ds.channelID, &discordgo.MessageEmbed{
			URL:         notification.Link,
			Title:       notification.Subject,
			Description: notification.CourseName,
			Footer:      &discordgo.MessageEmbedFooter{Text: notification.Date},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: notification.CourseImg,
			},
		})
	}
}
