package bot

import (
	"TwitchBot/internal/bot/commands"
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

const (
	common   string = "common"
	duel            = "duel"
	moderate        = "moderate"
)

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
	go func() {
		if strings.HasPrefix(message.Message, t.Prefix) { // cansel messages without prefix
			answer, commandType, err := t.findCommand(message.Message[len(t.Prefix):])
			if err != nil {
				fmt.Printf("error to find command: %s error: %s", message.Message, err.Error())
				return
			}

			if answer != "" {
				var mes string
				if commandType == common {
					mes, err = commands.CompileMessage(message, answer)
					if err != nil {
						fmt.Printf("error compile message: %s error: %s", answer, err.Error())
					}
					go t.sendMessage(mes)
				}

				if commandType == duel {
					mes, err = commands.CompileDuel(message, answer)
					if err != nil {
						fmt.Printf("error compile message: %s error: %s", answer, err.Error())
					}
					go t.sendMessage(mes)
					// TODO: Here should be goroutine call func that choose duel winner
				}

			}
		}
	}()
}

//findCommand start search for command in DB
func (t *userThread) findCommand(command string) (string, string, error) {
	var answer string
	var err error
	var find bool

	if t.isDuelEnabled() {
		answer, find, err = database.DB.FindDuelCommand(context.Background(), t.ChannelName)
		if err != nil {
			return "", "", err
		}
		if find {
			return answer, duel, nil
		}
	}
	if t.isCommonEnabled() {
		answer, find, err = database.DB.FindCommand(context.Background(), t.ChannelName, command)
		if err != nil {
			return "", "", err
		}
		if find {
			return answer, common, nil
		}
	}
	if t.isModerateEnabled() {
		return answer, moderate, nil
	}

	return answer, "", nil
}

func (t *userThread) isCommonEnabled() bool {
	for _, module := range t.Modules {
		if module == common {
			return true
		}
	}

	return false
}

func (t *userThread) isDuelEnabled() bool {
	for _, module := range t.Modules {
		if module == duel {
			return true
		}
	}

	return false
}

func (t *userThread) isModerateEnabled() bool {
	for _, module := range t.Modules {
		if module == moderate {
			return true
		}
	}

	return false
}

//sendMessage send message
func (t *userThread) sendMessage(message string) {
	t.Client.Say(t.ChannelName, message)
}
