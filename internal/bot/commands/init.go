package commands

import (
	"math/rand"
	"regexp"
	"time"
)

var reAuthorName *regexp.Regexp
var reChance *regexp.Regexp
var reRandomChatter *regexp.Regexp

var reDuelStatsAuthor *regexp.Regexp
var reDuelStatsOppo *regexp.Regexp
var reDuelOppo *regexp.Regexp
var reDuelWinner *regexp.Regexp
var reDuelLoser *regexp.Regexp

var reSetTitle *regexp.Regexp

//REInit initialize regular expression for special expression that should be replaced
// by some expression like number or message author name
func REInit() error {
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

	reDuelStatsAuthor, err = regexp.Compile(`\{%author_duel_stats%}`)
	if err != nil {
		return err
	}

	reDuelStatsOppo, err = regexp.Compile(`\{%oppo_duel_stats%}`)
	if err != nil {
		return err
	}

	reDuelOppo, err = regexp.Compile(`\{%oppo_duel%}`)
	if err != nil {
		return err
	}

	reDuelWinner, err = regexp.Compile(`\{%duel_winner%}`)
	if err != nil {
		return err
	}

	reDuelLoser, err = regexp.Compile(`\{%duel_loser%}`)
	if err != nil {
		return err
	}

	return nil
}
