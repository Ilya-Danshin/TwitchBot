package channel_interaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type userData struct {
	Id                string `json:"id"`
	Login             string `json:"login"`
	Display_name      string `json:"display_name"`
	TypeUser          string `json:"type"`
	Broadcaster_type  string `json:"broadcaster_type"`
	Description       string `json:"description"`
	Profile_image_url string `json:"profile_image_url"`
	Offline_image_url string `json:"offline_image_url"`
	View_count        int    `json:"view_count"`
	Email             string `json:"email"`
	Created_at        string `json:"created_at"`
}

type userInfo struct {
	data []userData
}

func (cc *ChatClient) SetTitle(channel, title string) error {
	// First step - get broadcaster(user) ID
	id, err := cc.getUserID(channel)
	if err != nil {
		return err
	}

	urlEnd, err := url.Parse("https://api.twitch.tv/helix/channels")
	if err != nil {
		return err
	}

	query := urlEnd.Query()
	query.Add("broadcaster_id", id)

	urlEnd.RawQuery = query.Encode()

	resp, err := cc.patch(options{
		endpoint: urlEnd.String(),
	}, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("missing or invalid parameter")
	}
	if resp.StatusCode == http.StatusInternalServerError {
		return fmt.Errorf("internal server error; failed to update channel")
	}
	return fmt.Errorf("unknown responce status")
}

func (cc *ChatClient) getUserInfo(username string) (*http.Response, error) {
	urlReq, err := url.Parse("https://api.twitch.tv/helix/users")
	if err != nil {
		return nil, err
	}

	query := urlReq.Query()
	query.Add("login", username)

	urlReq.RawQuery = query.Encode()

	resp, err := cc.get(options{
		endpoint: urlReq.String(),
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (cc *ChatClient) getUserID(username string) (string, error) {
	resp, err := cc.getUserInfo(username)
	if err != nil {
		return "", err
	}

	var info userInfo
	var body []byte
	_, err = resp.Body.Read(body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		return "", err
	}

	return info.data[0].Id, nil
}
