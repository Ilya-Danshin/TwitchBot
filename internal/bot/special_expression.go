package bot

import (
	"TwitchBot/internal/channel_interaction"
	"fmt"
	"github.com/gempir/go-twitch-irc/v3"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var reAuthorName *regexp.Regexp
var reChance *regexp.Regexp
var reRandomChatter *regexp.Regexp

//reInit initialize regular expression for special expression that should be replaced
// by some expression like number or message author name
func reInit() error {
	var err error

	rand.Seed(time.Now().Unix())

	reAuthorName, err = regexp.Compile(`\{%author_name%}`)
	if err != nil {
		return err
	}

	reChance, err = regexp.Compile(`\{%chance:-?\d+:-?\d+%}`)
	if err != nil {
		return err
	}

	reRandomChatter, err = regexp.Compile(`\{%random_chatter%}`)
	if err != nil {
		return err
	}

	return nil
}

//compileMessage start all find & replace function
func compileMessage(message twitch.PrivateMessage, answer string) (string, error) {
	mes, err := compileAuthorName(message.User.Name, answer)
	if err != nil {
		return "", err
	}

	mes, err = compileChance(mes)
	if err != nil {
		return "", err
	}

	mes, err = compileRandomChatter(message.Channel, answer)
	if err != nil {
		return "", err
	}

	return mes, nil
}

//compileAuthorName replace {%author_name%} by message author name
func compileAuthorName(author, message string) (string, error) {
	res := reAuthorName.FindAllString(message, -1)

	for _, expr := range res {
		message = strings.Replace(message, expr, author, 1)
	}

	return message, nil
}

//compileChance replace {%num1:num2%} by random integer in range (num1-num2)
func compileChance(message string) (string, error) {
	res := reChance.FindAllString(message, -1)

	for _, expr := range res {

		highBorder, lowBorder, err := getInterval(expr)
		if err != nil {
			return "", err
		}

		if highBorder < lowBorder {
			return "", fmt.Errorf("abs high border lower than abs low border")
		}

		result := rand.Intn(highBorder-lowBorder) + lowBorder

		message = strings.Replace(message, expr, strconv.Itoa(result), 1)
	}

	return message, nil
}

//getInterval from string {%num1:num2%} get low and bottom border
func getInterval(interval string) (int, int, error) {
	interval = interval[9 : len(interval)-2] // 9 = 2('{%') + 7('chance:')
	inter := strings.Split(interval, ":")

	lowBorder, err := strconv.Atoi(inter[0])
	if err != nil {
		return 0, 0, err
	}

	highBorder, err := strconv.Atoi(inter[1])
	if err != nil {
		return 0, 0, err
	}

	return highBorder, lowBorder, nil
}

func compileRandomChatter(channel, message string) (string, error) {
	res := reRandomChatter.FindAllString(message, -1)

	randomChatter, err := getRandomChatter(channel)
	if err != nil {
		return "", err
	}

	for _, expr := range res {
		message = strings.Replace(message, expr, randomChatter, 1)
	}

	return message, nil
}

func getRandomChatter(channel string) (string, error) {
	info, err := channel_interaction.GetChannelChatters(channel)
	if err != nil {
		return "", err
	}
	chatters := info.Chatters
	allChatters := append(chatters.Broadcaster, chatters.Vips...)
	allChatters = append(allChatters, chatters.Moderators...)
	allChatters = append(allChatters, chatters.Staff...)
	allChatters = append(allChatters, chatters.Admins...)
	allChatters = append(allChatters, chatters.GlobalMods...)
	allChatters = append(allChatters, chatters.Viewers...)

	randIndex := rand.Intn(len(allChatters))

	return allChatters[randIndex], nil
}
