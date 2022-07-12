package channel_interaction

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Chatters struct {
	Broadcaster []string
	Vips        []string
	Moderators  []string
	Staff       []string
	Admins      []string
	GlobalMods  []string
	Viewers     []string
}

type ChatInfo struct {
	Links        string
	ChatterCount int
	Chatters     Chatters
}

func (cc *ChatClient) GetChannelChatters(channel string) (*ChatInfo, error) {

	//resp, err := http.Get("http://tmi.twitch.tv/group/user/" + strings.ToLower(channel) + "/chatters")
	resp, err := cc.get(options{
		endpoint: "http://tmi.twitch.tv/group/user/" + strings.ToLower(channel) + "/chatters",
	})
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}

	var info ChatInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (cc *ChatClient) IsChatterInChat(channel, nickname string) (bool, error) {
	info, err := cc.GetChannelChatters(channel)
	if err != nil {
		return false, err
	}

	chatters := info.Chatters
	// In every group chatter in alphabetical order so try to find in largest group
	find := userSearch(chatters.Viewers, nickname)
	if find {
		return true, nil
	}
	find = userSearch(chatters.Vips, nickname)
	if find {
		return true, nil
	}
	find = userSearch(chatters.Moderators, nickname)
	if find {
		return true, nil
	}
	find = userSearch(chatters.Broadcaster, nickname)
	if find {
		return true, nil
	}
	find = userSearch(chatters.Staff, nickname)
	if find {
		return true, nil
	}
	find = userSearch(chatters.GlobalMods, nickname)
	if find {
		return true, nil
	}
	find = userSearch(chatters.Admins, nickname)
	if find {
		return true, nil
	}
	return false, nil
}

// In every group strings should be in alphabetical order, but sometimes his not, so can't use binary search
func userSearch(arr []string, user string) bool {
	for _, nickname := range arr {
		if nickname == user {
			return true
		}
	}

	return false
}
