package Dispatcher

import (
	"fmt"
	"github.com/LastSprint/feedback_bot/Steve/Models"
	"reflect"
)

type Handler interface {
	Handle(event Models.SlackEvent) error
}

type SlackEventDispatcher struct {
	handlers map[EventTypeKey]Handler
}

func (s *SlackEventDispatcher) RegisterHandler(key EventTypeKey, handler Handler) {
	if s.handlers == nil {
		s.handlers = map[EventTypeKey]Handler{}
	}
	s.handlers[key] = handler
}

func (s *SlackEventDispatcher) Dispatch(event Models.SlackEvent) error {
	key, err := StringToEventTypeKey(event.EventValue.Type)
	if err != nil {
		return fmt.Errorf("%s can't match this event %w", reflect.TypeOf(s), err)
	}

	handler, ok := s.handlers[key]

	if !ok {
		return fmt.Errorf("there wasn't any handler for key %s", key)
	}

	return handler.Handle(event)
}
