package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
	"strings"
	"time"
  "math/rand"
)

var soundFiles = []string{"Brann.dca", "Death.dca", "OhBrother.dca", "Secrets.dca"}

func init() {
	flag.StringVar(&token, "t", "", "Account Token")
	flag.Parse()
}

var token string

func main() {
	if token == "" {
		fmt.Println("No token provided. Please run: Brann -t <bot token>")
		return
	}

	err := loadSound()
	if err != nil {
		fmt.Println("Error loading sounds: ", err)
		fmt.Println("Please copy the sounds to this directory.")
		return
	}

	// Create a new Discord session using the provided token.
	dg, err := discordgo.New(token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(guildCreate)
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("Brann BotBeard is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	_ = s.UpdateStatus(0, "!!brann")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, "!!brann") {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			return
		}

		bufferID := 0

		switch m.Content {
		case "!!brann":
			bufferID = rand.Intn(4)
		case "!!brann trumpets":
			bufferID = 0
		case "!!brann death":
			bufferID = 1
		case "!!brann ohbrother":
			bufferID = 2
		case "!!brann secrets":
			bufferID = 3
		default:
			fmt.Println("Unknown Command: ", m.Content)
			return
		}
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playSound(s, g.ID, vs.ChannelID, bufferID)
				if err != nil {
					fmt.Println("Error playing sound:", err)
				}
				return
			}
		}
	}
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable != nil {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Brann is ready! Type !!brann while in a voice channel to play a sound.")
			return
		}
	}
}

var buffers = [][][]byte{}

func loadSound() error {

	for _, sound := range soundFiles {

		file, err := os.Open(sound)

		if err != nil {
			fmt.Println("Error opening dca file :", err)
			return err
		}

		var opuslen int16

		var buffer = make([][]byte, 0)

		for {
			err = binary.Read(file, binary.LittleEndian, &opuslen)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			if err != nil {
				fmt.Println("Error reading from dca file :", err)
				return err
			}
			InBuf := make([]byte, opuslen)
			err = binary.Read(file, binary.LittleEndian, &InBuf)
			if err != nil {
				fmt.Println("Error reading from dca file :", err)
				return err
			}
			buffer = append(buffer, InBuf)
		}
		buffers = append(buffers, buffer)
	}
	return nil
}

var isSoundPlaying = false

func playSound(s *discordgo.Session, guildID, channelID string, bufferID int) (err error) {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)

    isSoundPlaying = true

  	if err != nil {
  		return err
  	}
  	time.Sleep(250 * time.Millisecond)
  	_ = vc.Speaking(true)
  	for _, buff := range buffers[bufferID] {
  		vc.OpusSend <- buff
  	}
  	_ = vc.Speaking(false)
  	time.Sleep(250 * time.Millisecond)
  	_ = vc.Disconnect()
  	return nil
}
