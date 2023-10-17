package received

import (
	"context"

	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/pkg/client"
)

const (
	RECE_KEY_GPT = "@gpt "
)

type ReveResp interface {
	GptResp
}

type GptResp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int    `json:"Index"`
		FinishReason string `json:"finish_reason"`
		Message      struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}
	}
	Usage struct {
		PromptTokens     string `json:"prompt_tokens"`
		CompletionTokens string `json:"completion_tokens"`
		TotalTokens      string `json:"total_tokens"`
	}
}

type RequestInterface[T ReveResp] interface {
	Send(ctx context.Context, conf conf.ReveConfig, payload client.Payload) (T, error)
}

type EventInterface interface {
	Event(ctx context.Context, conf conf.Config, payload client.Payload) (string, error)
}
