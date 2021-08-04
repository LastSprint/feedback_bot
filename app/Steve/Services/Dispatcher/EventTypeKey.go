package Dispatcher

import "fmt"

type EventTypeKey string

const (
	EventTypeKeyMessage        EventTypeKey = "message"
	EventTypeKeyAddReaction    EventTypeKey = "reaction_added"
	EventTypeKeyRemoveReaction EventTypeKey = "reaction_removed"
)

func StringToEventTypeKey(val string) (EventTypeKey, error) {
	switch EventTypeKey(val) {
	case EventTypeKeyMessage:
		return EventTypeKeyMessage, nil
	case EventTypeKeyAddReaction:
		return EventTypeKeyAddReaction, nil
	case EventTypeKeyRemoveReaction:
		return EventTypeKeyRemoveReaction, nil
	}

	return "", fmt.Errorf("%s didn't match possible EventTypeKey values")
}
