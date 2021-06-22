package Repo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// SlackRepo contains methods for calling slack API
type SlackRepo struct {
	// AuthToken is token for API access
	// you can use both (user and bot)
	AuthToken string
}

// PostMessageToChat will post one message to specific channel in declared thread
// If threadId was specified then message will be posted in thread (in specific channel)
// otherwise - will be posted right in channel
func (r *SlackRepo) PostMessageToChat(message, channel, threadId string) error {

	msg := struct {
		Channel  string `json:"channel"`
		Text     string `json:"text"`
		ThreadTs string `json:"thread_ts"`
	}{
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

	headers.Add("Authorization", "Bearer "+r.AuthToken)
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

	defer func() {
		if err = resp.Body.Close(); err != nil {
			fmt.Printf("[WARN] Couldn't close response body for %v", resp.Request.URL)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("POST on %s returned %v status code with body %s", req.URL, resp.StatusCode, string(data)))
	}

	return nil
}
