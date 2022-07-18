package bot

import (
	"fmt"

	"TwitchBot/internal/bot/commands"

	"github.com/gempir/go-twitch-irc/v3"
)

func (t *channelThread) isCommonEnabled() bool {
	return isContain(common, t.Modules)
}

//commonCommandHandler handler for all common chat commands
func (t *channelThread) commonCommandHandler(message twitch.PrivateMessage, answer string) {
	mes, err := commands.CompileCommonMessage(message, answer)
	if err != nil {
		fmt.Printf("error compile message: %s error: %s", answer, err.Error())
		return
	}
	go t.sendMessage(mes)
}
