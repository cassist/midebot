package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/jasonlvhit/gocron"
)

var (
	// discordgo session
	discord *discordgo.Session

	// OWNER
	OWNER string

	//Token
	authToken string

	cactpotTime time.Time

	// fresh milk guild and channel IDs
	freshMilk = "237317207217012737"
	glamtrash = "238798619858173963"
	sss       = "238805318577029121"

	//autotranslate emojis
	leftAT  = "<:left:240181451876859904>"
	rightAT = "<:right:240181451885117440>"

	//other emojis
	pray   = "<:pray:237783447895015424>"
	mogeko = "<:mogeko:240183110422102016>"
)

// func handleControlCommand(s *discordgo.Session, m *discordgo.MessageCreate, msg []string, guild *discordgo.Guild) {
// 	// insert code for administrative commands
// 	return
// }

func handleSlash(s *discordgo.Session, m *discordgo.MessageCreate, msg []string, guild *discordgo.Guild) {
	user := "<@!" + m.Author.ID + ">"
	switch msg[0] {
	case "/at":
		newMsg := leftAT + msg[1] + rightAT
		sendMessage(s, m.ChannelID, newMsg, "autotranslated")

		// until discord decides to add permissions to edit others' messages, this cannot be done
		// _, err := s.ChannelMessageEdit(m.ChannelID, m.ID, newMsg)
		// if err != nil {
		// 	log.WithFields(log.Fields{
		// 		"channel": m.ChannelID,
		// 		"message": m.ID,
		// 		"err":     err,
		// 	}).Warning("Failed to autotranslate message")
		//}
	case "/aqua":
		var newMsg string
		watdo := randomNumberGenerator(4)
		switch watdo {
		case 0:
			newMsg = "Let's go **right**! " + pray
		case 1:
			newMsg = "Let's go **left**! " + pray
		case 2:
			newMsg = mogeko + " **Right** side, もげげ。"
		case 3:
			newMsg = mogeko + " Go **left**, もげ!"
		}
		sendMessage(s, m.ChannelID, newMsg, "aqua")
	case "/cactpot":
		newMsg := "Cactpot draws are Saturdays at 9pm EST.\r\n"
		newMsg += "My cactpot reminder is set to go off at " + cactpotTime.Format("Mon Jan 2 15:04:05") + "."
		sendMessage(s, m.ChannelID, newMsg, "cactpot")
	case "/draw":
		draw(s, m.ChannelID, m.Author.ID)
	case "/info":
		newMsg := "My commands: **/aqua**: Aquapolis helper. \r\n"
		newMsg += "**/at {message}**: I'll reply with " + leftAT + "{message}" + rightAT + ".\r\n"
		newMsg += "**/cactpot**: I send out a weekly reminder about Jumbo Cactpot. Use this to confirm the reminder time.\r\n"
		newMsg += "**/draw**: Use the Astrologian skill " + leftAT + "Draw" + rightAT + ".\r\n"
		newMsg += "**/random**: Generates a random number like the FF command.\r\n\r\n"
		newMsg += "I automatically pin image uploads in the screenshot and glam channels!\r\n"
		newMsg += "Midebot is written in golang using discordgo, by Cassis Milk of Jenova"
		sendMessage(s, m.ChannelID, newMsg, "info")
	case "/pray":
		newMsg := user + " prays solemnly. " + pray
		sendMessage(s, m.ChannelID, newMsg, "pray")
	case "/random":
		newMsg := "Random! " + user + " rolls a 🎲" + fmt.Sprintf("%v", randomNumberGenerator(1000)) + "."
		sendMessage(s, m.ChannelID, newMsg, "rng")
	}
}

func sendMessage(s *discordgo.Session, cid string, msg string, msgType string) {
	_, err := s.ChannelMessageSend(cid, msg)
	if err != nil {
		log.WithFields(log.Fields{
			"channel": cid,
			"message": msg,
			"err":     err,
		}).Warning("Failed to send" + msgType + " message")
	}
}

func pinSS(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, attach := range m.Attachments {
		if isImage(attach.Filename) {
			err := s.ChannelMessagePin(m.ChannelID, m.ID)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Warning("Failed to pin message")
			}
			return
		}
	}
}

func isImage(file string) bool {
	re := regexp.MustCompile(".*[.].{3,4}")
	if !re.MatchString(file) {
		return false
	}
	re = regexp.MustCompile("[.]")
	filetype := re.Split(file, -1)
	switch filetype[1] {
	case `jpg`, `jpeg`, `png`, `gif`, `bmp`:
		return true
	default:
		return false
	}
}

func randomNumberGenerator(n int) (num int) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	num = r.Intn(n)
	return num
}

func draw(s *discordgo.Session, cid string, author string) {
	newMsg := "<@!" + author + "> has drawn a"
	card := randomNumberGenerator(100)
	switch {
	case card <= 9:
		//balance
		newMsg += " **Balance** card!"
		sendMessage(s, cid, newMsg, "balance")
		sendCard(s, "balance.png", cid)
	case (card > 9 && card <= 34):
		//bole
		newMsg += " **Bole** card!"
		sendMessage(s, cid, newMsg, "bole")
		sendCard(s, "bole.png", cid)
	case (card > 34 && card <= 49):
		//arrow
		newMsg += "n **Arrow** card!"
		sendMessage(s, cid, newMsg, "arrow")
		sendCard(s, "arrow.png", cid)
	case (card > 49 && card <= 64):
		//spear
		newMsg += " **Spear** card!"
		sendMessage(s, cid, newMsg, "spear")
		sendCard(s, "spear.png", cid)
	case (card > 64 && card <= 79):
		//ewer
		newMsg += " **Ewer** card!"
		sendMessage(s, cid, newMsg, "ewer")
		sendCard(s, "ewer.png", cid)
	case (card > 79 && card <= 99):
		//spire
		newMsg += " **Spire** card!"
		sendMessage(s, cid, newMsg, "spire")
		sendCard(s, "spire.png", cid)
	}
	return
}

func sendCard(s *discordgo.Session, name string, cid string) {
	if name != "" {
		f, err := os.Open(fmt.Sprintf("ast/%s", name))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Warning("Failed to upload card.")
		}
		defer f.Close()
		s.ChannelFileSend(cid, name, f)
		return
	}
}

func onGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Guild.Unavailable != nil {
		return
	}
	log.Printf("joined guild: %s <%s>", g.Name, g.ID)

	for _, channel := range g.Guild.Channels {
		if channel.ID == g.Guild.ID {
			s.ChannelMessageSend(channel.ID, "<pant><pant>")
			return
		}
	}
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// check ss to pin first, since it won't have command syntax
	if (m.Attachments != nil || m.Embeds != nil) && (m.ChannelID == sss || m.ChannelID == glamtrash) {
		pinSS(s, m)
	}

	// eliminate messages that do not match command syntax
	if len(m.Content) <= 0 || (m.Content[0] != '/' && m.Content[0] != '!') {
		return
	}

	//msg := m.ContentWithMentionsReplaced()        // this is causing issues with auto-translating names
	msg := m.Content
	parts := strings.SplitN(msg, " ", 2)

	channel, _ := discord.State.Channel(m.ChannelID)
	if channel == nil {
		log.WithFields(log.Fields{
			"channel": m.ChannelID,
			"message": m.ID,
		}).Warning("Failed to grab channel")
		return
	}

	guild, _ := discord.State.Guild(channel.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild":   channel.GuildID,
			"channel": channel,
			"message": m.ID,
		}).Warning("Failed to grab guild")
		return
	}

	if m.Content[0] == '/' {
		handleSlash(s, m, parts, guild)
	}

	// if m.Content[0] == '!' {

	// }

	// If this is a mention, it should come from the owner (otherwise we don't care)
	// if len(m.Mentions) > 0 && m.Author.ID == OWNER && len(parts) > 0 {
	// 	mentioned := false
	// 	for _, mention := range m.Mentions {
	// 		mentioned = (mention.ID == s.State.Ready.User.ID)
	// 		if mentioned {
	// 			break
	// 		}
	// 	}

	// 	if mentioned {
	// 		handleControlCommand(s, m, parts, guild)
	// 	}
	// 	return
	// }
	return
}

func sendCactpot(s *discordgo.Session) {
	newMsg := "@everyone Reminder! Jumbo Cactpot drawings are at 9PM!"
	_, err := s.ChannelMessageSend(freshMilk, newMsg)
	if err != nil {
		log.WithFields(log.Fields{
			"channel": freshMilk,
			"message": newMsg,
			"err":     err,
		}).Warning("Failed to send cactpot message")
	}
}

func jumboReminder(s *discordgo.Session) {
	for {
		gocron.Every(1).Saturday().At("20:55").Do(sendCactpot, s)
		_, cactpotTime = gocron.NextRun()
		fmt.Println(cactpotTime)

		<-gocron.Start()
	}
	// me making stupid decisions

	// var timeNow time.Time
	// for {
	// 	timeNow = time.Now()
	// 	if timeNow.Weekday() == time.Saturday {
	// 		if timeNow.Hour() == 21 {
	// 			newMsg := "@everyone Reminder! Jumbo Cactpot drawings are at 10PM!"
	// 			_, err := s.ChannelMessageSend(freshMilk, newMsg)
	// 			if err != nil {
	// 				log.WithFields(log.Fields{
	// 					"channel": freshMilk,
	// 					"message": newMsg,
	// 					"err":     err,
	// 				}).Warning("Failed to send cactpot message")
	// 			}
	// 		}
	// 	} else {
	// 		time.Sleep(time.Hour)
	// 	}
	// }
}

func main() {

	var (
		Token    = flag.String("t", "", "Discord Authentication Token")
		Email    = flag.String("e", "", "Discord Email")
		Password = flag.String("p", "", "Discord Password")
		Owner    = flag.String("o", "", "Owner ID")
		err      error
	)
	flag.Parse()

	if *Owner != "" {
		OWNER = *Owner
	}
	if *Token != "" {
		authToken = *Token
	}

	// Create a discord session
	log.Info("Starting discord session...")
	if *Email != "" && *Password != "" {
		discord, err = discordgo.New(*Email, *Password, *Token)
	} else {
		discord, err = discordgo.New(authToken)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create discord session")
		return
	}

	discord.AddHandler(onGuildCreate)
	discord.AddHandler(onMessageCreate)

	err = discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to open discord websocket connection")
		return
	}

	// We're running!
	log.Info("Mide is ready to go forward and back")

	go jumboReminder(discord)

	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

}
