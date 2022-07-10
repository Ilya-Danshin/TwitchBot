package bot

import (
	"github.com/gempir/go-twitch-irc/v3"
)

//moderateCommandHandler handler for all moderate chat commands
func (t *channelThread) moderateCommandHandler(message twitch.PrivateMessage, answer string) {
	return
}

func (t *channelThread) isModerateEnabled() bool {
	for _, module := range t.Modules {
		if module == moderate {
			return true
		}
	}

	return false
}
