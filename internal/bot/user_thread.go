package bot

import (
	"TwitchBot/internal/bot/commands"
	"context"
	"fmt"
	"strings"
	"time"

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
	Duel        string
	DuelDelay   int
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
		Duel:        user.Duel,
		DuelDelay:   user.DuelDelay,
		Modules:     user.Modules,
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
	if strings.HasPrefix(message.Message, t.Prefix) { // cansel messages without prefix
		answer, commandType, err := t.findCommand(message.Message[len(t.Prefix):])
		if err != nil {
			fmt.Printf("error to find command: %s error: %s", message.Message, err.Error())
			return
		}

		if answer != "" {
			if commandType == common {
				t.commonCommandHandler(message, answer)
				return
			}
			if commandType == duel {
				t.duelCommandHandler(message, answer)
				return
			}
			if commandType == moderate {
				t.moderateCommandHandler(message, answer)
				return
			}

		}
	}

}

//commonCommandHandler handler for all common chat commands
func (t *userThread) commonCommandHandler(message twitch.PrivateMessage, answer string) {
	mes, err := commands.CompileMessage(message, answer)
	if err != nil {
		fmt.Printf("error compile message: %s error: %s", answer, err.Error())
	}
	go t.sendMessage(mes)
}

//duelCommandHandler handler for duel chat command
func (t *userThread) duelCommandHandler(message twitch.PrivateMessage, answer string) {
	mes, oppo, err := commands.CompileDuel(message, answer, t.Prefix, t.Duel)
	if err != nil {
		fmt.Printf("error compile duel message: %s error: %s", answer, err.Error())
	}
	if mes != "" { // If there is empty message than duel was canceled or was error
		go t.sendMessage(mes)
		// Async call for
		go func() {
			mes, err = commands.GetDuelWinner(message, oppo)
			if err != nil {
				fmt.Printf("error get duel winner message: %s error: %s", answer, err.Error())
			}
			time.Sleep(time.Second * time.Duration(t.DuelDelay))

			if mes != "" {
				go t.sendMessage(mes)
			}
		}()
	}
}

//moderateCommandHandler handler for all moderate chat commands
func (t *userThread) moderateCommandHandler(message twitch.PrivateMessage, answer string) {
	return
}

//findCommand start search for command in DB
func (t *userThread) findCommand(command string) (string, string, error) {
	var answer string
	var err error
	var find bool

	if t.isDuelEnabled() {
		if strings.HasPrefix(command, t.Duel) {
			answer, find, err = database.DB.FindDuelCommand(context.Background(), t.ChannelName, t.Duel)
			if err != nil {
				return "", "", err
			}
			if find {
				return answer, duel, nil
			}
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
