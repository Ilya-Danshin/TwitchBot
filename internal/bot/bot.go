package bot

import (
	"TwitchBot/internal/bot/commands"
	"fmt"

	"TwitchBot/config"
)

var errorsChan chan error
var threads []*userThread

// InitBot just run all bot process
func InitBot(users []*config.User, botSettings *config.BotSettings) error {
	err := commands.REInit()
	if err != nil {
		return err
	}
	errorsChan = make(chan error)
	for _, user := range users {
		threads = append(threads, NewUserThread(user, botSettings))
	}

	return nil
}

// LoopBot listening to channel, where bot send errors
func LoopBot() {
	for _, thread := range threads {
		go thread.Run()
	}

	for {
		select {
		case err := <-errorsChan:
			fmt.Printf("error: %s", err.Error())
		}
	}
}
