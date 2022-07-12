package bot

import (
	"TwitchBot/database"
	"TwitchBot/internal/bot/commands"
	"context"
	"fmt"
	"github.com/gempir/go-twitch-irc/v3"
	"strings"
)

func (t *channelThread) isSubdayEnabled() bool {
	for _, module := range t.Modules {
		if module == subday {
			return true
		}
	}

	return false
}

func (t *channelThread) subdayCommandHandler(message twitch.PrivateMessage, answer string) {
	// TODO: Add sub/follower check with settings from JSON
	split := strings.SplitN(message.Message, " ", 2)
	order := split[1]
	err := database.DB.AddNewSubdayOrder(context.Background(), message.Channel, message.User.Name, order)
	if err != nil {
		fmt.Printf("error add new subday order error: %s", err.Error())
		return
	}

	mes, err := commands.CompileSubdayMessage(message, answer)
	if err != nil {
		fmt.Printf("error compile subday message error: %s", err.Error())
		return
	}

	if mes != "" {
		t.sendMessage(mes)
	}
}
