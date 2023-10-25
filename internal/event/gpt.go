package event

import (
	"context"
	"fmt"
	"os/exec"

	"sync"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/jsonx"

	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
	openai "github.com/sashabaranov/go-openai"
)

var lock sync.Mutex
var lockFlag bool

type gpt struct {
	keywords []string
}

type GptReqMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (e *gpt) Event(ctx context.Context, svcCtx *svc.ServiceContext, payload *client.Payload) (string, error) {
	qust := payload.String()
	idx := keywords(qust, []string{svcCtx.Config.GPT.Keywords, "justyolo "})
	if idx < 0 {
		return "", nil
	}
	if lockFlag {
		return "", errors.New("忙着呢，稍候在问!!")
	}
	lock.Lock()
	lockFlag = true

	defer func() {
		lockFlag = false
		lock.Unlock()
	}()

	qust = qust[idx:]
	if len(qust) == 0 {
		return "", errors.New("请说出你想问的问题")
	}
	return e.proxy(ctx, svcCtx, qust)
	// return e.qustion(ctx, svcCtx, qust)
}

func (e *gpt) qustion(ctx context.Context, svcCtx *svc.ServiceContext, qust string) (string, error) {
	fmt.Printf("开始提问:%v", qust)
	cmd := exec.Command("curl", "--max-time", "300", "--request", "POST", `https://api.f2gpt.com/v1`,
		"-H", "Content-Type: application/json",
		"-H", fmt.Sprintf("Authorization: Bearer %v", svcCtx.Config.GPT.Key),
		"-d", fmt.Sprintf(`{"model":"%v","messages":[{"role":"%v","content":"%v"}]}`,
			openai.GPT3Dot5Turbo, openai.ChatMessageRoleUser, qust))
	out, err := cmd.Output()

	resp := &openai.ChatCompletionResponse{}
	err = jsonx.Unmarshal(out, resp)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	chLen := len(resp.Choices)
	if chLen == 0 {
		fmt.Println(string(out))
		return "", fmt.Errorf("ChatCompletion len = 0 or error: %v", err)
	}

	return resp.Choices[chLen-1].Message.Content, nil
}

func (e *gpt) proxy(ctx context.Context, svcCtx *svc.ServiceContext, qust string) (string, error) {
	config := openai.DefaultConfig(svcCtx.Config.GPT.Key)
	config.BaseURL = "https://api.f2gpt.com/v1"
	fmt.Printf("开始提问%v", qust)
	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{{
				Role:    openai.ChatMessageRoleUser,
				Content: qust},
			},
		},
	)
	fmt.Println("返回", err, resp)
	chLen := len(resp.Choices)
	if err != nil || chLen == 0 {
		return "", fmt.Errorf("ChatCompletion len = 0 or error: %v", err)
	}

	return resp.Choices[chLen-1].Message.Content, nil
}
