package bot

import (
	"fmt"
	"time"

	"TwitchBot/internal/bot/commands"

	"github.com/gempir/go-twitch-irc/v3"
)

func (t *channelThread) isDuelEnabled() bool {
	for _, module := range t.Modules {
		if module == duel {
			return true
		}
	}

	return false
}

//duelCommandHandler handler for duel chat command
func (t *channelThread) duelCommandHandler(message twitch.PrivateMessage, answer string) {
	mes, oppo, err := commands.CompileDuelMessage(message, answer, t.Prefix, t.Duel.DuelWord)
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
			time.Sleep(time.Second * time.Duration(t.Duel.DuelDelay))

			if mes != "" {
				go t.sendMessage(mes)
			}
		}()
	}
}
