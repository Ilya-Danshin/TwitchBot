package commands

import (
	"TwitchBot/database"
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

func compileDuelNames(message, authorName, oppoName string) (string, error) {
	res := reDuelOppo.FindAllString(message, -1)

	for _, expr := range res {
		message = strings.Replace(message, expr, oppoName, 1)
	}

	return message, nil
}

func CompileDuel(message twitch.PrivateMessage, answer string) (string, error) {
	mes, err := compileAuthorName(answer, message.User.Name)
	if err != nil {
		return "", err
	}

	mes, err = compileChance(mes)
	if err != nil {
		return "", err
	}
	var oppo string
	oppo, err = getRandomChatter(message.Channel)
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

	mes, err = compileDuelNames(mes, message.User.Name, oppo)
	if err != nil {
		return "", err
	}

	mes, err = compileDuelStats(mes, authorStats, oppoStats)
	if err != nil {
		return "", err
	}

	return mes, nil
}
