package commands

import (
	"TwitchBot/database"
	"TwitchBot/internal/channel_interaction"
	"context"
	"github.com/gempir/go-twitch-irc/v3"
	"strings"
)

func compileDuelStats(message string, authorStats, oppoStats *database.DuelStats) (string, error) {
	res := reDuelStatsAuthor.FindAllString(message, -1)

	authorStatsString := authorStats.String()

	for _, expr := range res {
		message = strings.Replace(message, expr, authorStatsString, 1)
	}

	res = reDuelStatsOppo.FindAllString(message, -1)

	oppoStatsString := oppoStats.String()

	for _, expr := range res {
		message = strings.Replace(message, expr, oppoStatsString, 1)
	}

	return message, nil
}

func compileDuelNames(message, oppoName string) (string, error) {
	res := reDuelOppo.FindAllString(message, -1)

	for _, expr := range res {
		message = strings.Replace(message, expr, oppoName, 1)
	}

	return message, nil
}

func CompileDuel(message twitch.PrivateMessage, answer, prefix, duelCommand string) (string, error) {
	mes, err := compileAuthorName(answer, message.User.Name)
	if err != nil {
		return "", err
	}

	mes, err = compileChance(mes)
	if err != nil {
		return "", err
	}

	oppo, err := chooseDuelTarget(message, prefix, duelCommand)
	if err != nil {
		return "", err
	}

	authorStats, err := database.DB.FindDuelUser(context.Background(), message.User.Name)
	if err != nil {
		return "", err
	}
	oppoStats, err := database.DB.FindDuelUser(context.Background(), oppo)
	if err != nil {
		return "", err
	}

	mes, err = compileDuelNames(mes, oppo)
	if err != nil {
		return "", err
	}

	mes, err = compileDuelStats(mes, authorStats, oppoStats)
	if err != nil {
		return "", err
	}

	return mes, nil
}

func chooseDuelTarget(message twitch.PrivateMessage, prefix, duelCommand string) (string, error) {
	var oppo string
	var err error

	if len([]rune(message.Message)) > len([]rune(prefix+duelCommand)) {
		// If its target duel
		// Trim duel-word and '@'
		oppo = strings.TrimPrefix(message.Message, prefix+duelCommand+" ")
		if oppo[0] == '@' {
			oppo = strings.TrimPrefix(oppo, "@")
		}
		oppo = strings.ToLower(oppo)
		inChat, err := channel_interaction.IsChatterInChat(message.Channel, oppo)
		if err != nil {
			return "", err
		}
		if !inChat {
			return "", nil
		}
	} else {
		// If its non-target duel
		oppo, err = getRandomChatter(message.Channel)
		if err != nil {
			return "", err
		}
	}

	return oppo, nil
}
