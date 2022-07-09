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

func GetChannelChatters(channel string) (*ChatInfo, error) {
	resp, err := http.Get("http://tmi.twitch.tv/group/user/" + strings.ToLower(channel) + "/chatters")
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
