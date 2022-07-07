package bot

import (
	"TwitchBot/config"
	"fmt"
)

var errorsChan chan error
var threads []*userThread

// InitBot just run all bot process
func InitBot(users []*config.User, botSettings *config.BotSettings) {
	errorsChan = make(chan error)
	for _, user := range users {
		threads = append(threads, NewUserThread(user, botSettings))
	}
}

// LoopBot listening to channel, where bot send errors
func LoopBot() {
	for i, thread := range threads {
		go thread.Run(i)
	}

	for {
		select {
		case err := <-errorsChan:
			fmt.Printf("error: %s", err.Error())
		}
	}
}
