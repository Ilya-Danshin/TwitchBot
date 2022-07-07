package bot

import (
	"fmt"
	"github.com/gempir/go-twitch-irc/v3"

	"TwitchBot/config"
)

type bot struct {
	Nickname string
	Oauth    string
}

type userThread struct {
	BotSettings bot
	ErrorChan   chan error
	ChannelName string
	Modules     []string
}

//NewUserThread create new thread object
func NewUserThread(user *config.User, botCfg *config.BotSettings) *userThread {
	return &userThread{
		BotSettings: bot{
			Nickname: botCfg.Nickname,
			Oauth:    botCfg.Oauth,
		},
		ErrorChan:   errorsChan,
		ChannelName: user.Name,
		Modules:     user.Modules, // TODO: Create function to translate string name module to pointer to function
	}
}

//Run thread
func (t *userThread) Run(i int) {
	client := twitch.NewClient(t.BotSettings.Nickname, t.BotSettings.Oauth)

	client.OnPrivateMessage(t.MessageFilter)

	client.Join(t.ChannelName)

	err := client.Connect()
	t.ErrorChan <- err

	return
}

//MessageFilter do all work with received message
func (t *userThread) MessageFilter(message twitch.PrivateMessage) {
	fmt.Printf("Channel: %s Author: %s Message: %s\n", t.ChannelName, message.User.Name, message.Message)
}
