package bot

import (
	"fmt"

	"TwitchBot/internal/bot/commands"

	"github.com/gempir/go-twitch-irc/v3"
)

func (t *channelThread) isCommonEnabled() bool {
	for _, module := range t.Modules {
		if module == common {
			return true
		}
	}

	return false
}

//commonCommandHandler handler for all common chat commands
func (t *channelThread) commonCommandHandler(message twitch.PrivateMessage, answer string) {
	mes, err := commands.CompileCommonMessage(message, answer)
	if err != nil {
		fmt.Printf("error compile message: %s error: %s", answer, err.Error())
	}
	go t.sendMessage(mes)
}
