package bot

import (
	"TwitchBot/internal/channel_interaction"
	"github.com/gempir/go-twitch-irc/v3"
	"strings"
)

func (t *channelThread) isModerateEnabled() bool {
	for _, module := range t.Modules {
		if module == moderate {
			return true
		}
	}

	return false
}

//moderateCommandHandler handler for all moderate chat commands
func (t *channelThread) moderateCommandHandler(message twitch.PrivateMessage, answer string) {
	split := strings.SplitN(message.Message, " ", 2)
	command := answer
	title := split[1]

	switch command {
	case "settitle":
		client := channel_interaction.NewClient()
		client.SetTitle(message.Channel, title)
		return
	}

}
