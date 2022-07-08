package bot

import (
	"context"
	"fmt"
	"strings"

	"TwitchBot/config"
	"TwitchBot/database"

	"github.com/gempir/go-twitch-irc/v3"
)

type bot struct {
	Nickname string
	Oauth    string
}

type userThread struct {
	BotSettings bot
	ErrorChan   chan error
	ChannelName string
	Prefix      string
	Modules     []string

	Client *twitch.Client
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
		Prefix:      user.Prefix,
		Modules:     user.Modules, // TODO: Create function to translate string name module to pointer to function
	}
}

//Run thread
func (t *userThread) Run() {
	client := twitch.NewClient(t.BotSettings.Nickname, t.BotSettings.Oauth)

	client.OnPrivateMessage(t.messageFilter)

	client.Join(t.ChannelName)

	t.Client = client

	err := client.Connect()
	t.ErrorChan <- err

	return
}

//messageFilter do all work with received message
func (t *userThread) messageFilter(message twitch.PrivateMessage) {
	//fmt.Printf("Channel: %s Author: %s Message: %s\n", t.ChannelName, message.User.Name, message.Message)
	if strings.HasPrefix(message.Message, t.Prefix) { // cansel messages without prefix
		answer, err := t.findCommand(message.Message[len(t.Prefix):])
		if err != nil {
			fmt.Printf("error to find command: %s error: %s", message.Message, err.Error())
			return
		}
		if answer != "" {
			mes, err := compileMessage(message, answer)
			if err != nil {
				fmt.Printf("error compile message: %s error: %s", answer, err.Error())
			}
			t.sendMessage(mes)
		}
	}
}

func (t *userThread) findCommand(command string) (string, error) {
	answer, err := database.DB.FindCommand(context.Background(), t.ChannelName, command)
	if err != nil {
		return "", nil
	}
	return answer, nil
}

func (t *userThread) sendMessage(message string) {
	t.Client.Say(t.ChannelName, message)
}
