package bot

import (
	"context"
	"fmt"
	"strings"

	"TwitchBot/config"
	"TwitchBot/database"

	"github.com/gempir/go-twitch-irc/v3"
)

type duelSettings struct {
	DuelWord  string
	DuelDelay int
}

type channelThread struct {
	BotSettings bot
	ErrorChan   chan error
	ChannelName string
	Prefix      string
	Duel        duelSettings
	SubdayWord  string

	Modules []string

	Client *twitch.Client
	ctx    context.Context
}

const (
	common   string = "common"
	duel            = "duel"
	moderate        = "moderate"
	subday          = "subday"
)

//NewUserThread create new thread object
func NewUserThread(user *config.User, botCfg *config.BotSettings) *channelThread {
	return &channelThread{
		BotSettings: bot{
			Nickname: botCfg.Nickname,
			Oauth:    botCfg.Oauth,
		},
		ErrorChan:   errorsChan,
		ChannelName: user.Name,
		Prefix:      user.Prefix,
		Duel: duelSettings{
			DuelWord:  user.Duel,
			DuelDelay: user.DuelDelay,
		},
		SubdayWord: user.Subday,
		Modules:    user.Modules,
		ctx:        context.Background(),
	}
}

//Run thread
func (t *channelThread) Run() {
	client := twitch.NewClient(t.BotSettings.Nickname, t.BotSettings.Oauth)

	client.OnPrivateMessage(t.messageFilter)

	client.Join(t.ChannelName)

	t.Client = client

	err := client.Connect()
	t.ErrorChan <- err

	return
}

//messageFilter do all work with received message
func (t *channelThread) messageFilter(message twitch.PrivateMessage) {
	if !strings.HasPrefix(message.Message, t.Prefix) { // cansel messages without prefix
		return
	}

	answer, commandType, err := t.findCommand(message.Message[len(t.Prefix):])
	if err != nil {
		fmt.Printf("error to find command: %s error: %s", message.Message, err.Error())
		return
	}

	if answer == "" {
		return
	}

	switch commandType {
	case common:
		t.commonCommandHandler(message, answer)
	case duel:
		t.duelCommandHandler(message, answer)
	case moderate:
		t.moderateCommandHandler(message, answer)
	case subday:
		t.subdayCommandHandler(message, answer)
	}
}

//findCommand start search for command in DB
func (t *channelThread) findCommand(command string) (string, string, error) {
	var answer string
	var err error
	var find bool

	if strings.HasPrefix(command, t.Duel.DuelWord) {
		if t.isDuelEnabled() {
			answer, find, err = database.DB.FindDuelCommand(t.ctx, t.ChannelName)
			if err != nil {
				return "", "", err
			}
			if find {
				return answer, duel, nil
			}
		}
		return "", duel, nil
	}
	if strings.HasPrefix(command, t.SubdayWord) {
		if t.isSubdayEnabled() {
			answer, find, err = database.DB.FindSubdayCommand(t.ctx, t.ChannelName)
			if err != nil {
				return "", "", err
			}

			if find {
				return answer, subday, nil
			}

			return "", subday, nil
		}
	}
	if t.isCommonEnabled() {
		answer, find, err = database.DB.FindCommonCommand(t.ctx, t.ChannelName, command)
		if err != nil {
			return "", "", err
		}
		if find {
			return answer, common, nil
		}
	}
	if t.isModerateEnabled() {
		split := strings.SplitN(command, " ", 2)

		answer, find, err = database.DB.FindModerateCommand(t.ctx, t.ChannelName, split[0])
		if err != nil {
			return "", "", err
		}
		if find {
			return answer, moderate, nil
		}
	}

	return answer, "", nil
}

//sendMessage send message
func (t *channelThread) sendMessage(message string) {
	t.Client.Say(t.ChannelName, message)
}
