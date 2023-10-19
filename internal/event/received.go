package event

import (
	"context"
	"fmt"

	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
	"github.com/sashabaranov/go-openai"
)

type EventType string

const (
	EVENT_KEY_GPT EventType = "gpt"
)

type ReveResp interface {
	openai.ChatCompletionResponse
}

type EventInterface interface {
	Event(ctx context.Context, svcCtx *svc.ServiceContext, payload client.Payload) (string, error)
}

var eventHub = map[EventType]EventInterface{
	EVENT_KEY_GPT: newGpt(),
}

func MustEventFactory(ek EventType) EventInterface {
	event, ok := eventHub[ek]
	if !ok {
		panic(fmt.Sprintf("not defined event %v", ek))
	}
	return event
}
