package bot

import (
	"github.com/gempir/go-twitch-irc/v3"
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
	//split := strings.SplitN(message.Message, " ", 2)
	command := answer
	//title := split[1]

	switch command {
	case "settitle":

		return
	}

}
