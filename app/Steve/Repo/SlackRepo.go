package Repo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type SlackRepo struct {
	AuthToken string
}

func (r *SlackRepo) PostMessageToChat(message, channel, threadId string) error {
	type slack_message_body struct {
		Channel  string `json:"channel"`
		Text     string `json:"text"`
		ThreadTs string `json:"thread_ts"`
	}

	msg := slack_message_body{
		Channel:  channel,
		Text:     message,
		ThreadTs: threadId,
	}

	data, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/chat.postMessage", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	headers := req.Header

	headers.Add("Authorization", "Bearer " + r.AuthToken)
	headers.Add("Content-Type", "application/json")
	req.Header = headers

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	data, err = json.Marshal(resp.Body)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("POST on %s returned %v status code with body %s", req.URL, resp.StatusCode, string(data)))
	}

	return nil
}