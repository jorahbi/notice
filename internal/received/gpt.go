package received

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/jsonx"

	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/pkg/client"
	openai "github.com/sashabaranov/go-openai"
)

var lock sync.Mutex
var lockFlag bool

const GPT_URL = "https://api.openai.com/v1/chat/completions"

type Gpt struct {
	Model string `json:"model"`
	// Temperature float32      `json:"temperature"`
	Messages []*GptReqMsg `json:"messages"`
}

type GptReqMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewGpt() *Gpt {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:7890")
	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:7890")
	return &Gpt{}
}

func (gpt *Gpt) Event(ctx context.Context, svcConf conf.Config, payload client.Payload) (string, error) {
	if lockFlag {
		return "", errors.New("忙着呢，稍候在问!!")
	}
	lock.Lock()
	lockFlag = true

	defer func() {
		lockFlag = false
		lock.Unlock()
	}()
	qust := payload.String()
	idx := strings.Index(qust, svcConf.GptKeywords)
	fmt.Println(qust, idx)
	if idx < 0 {
		return "", nil
	}
	qust = strings.Trim(qust[len(svcConf.GptKeywords):], " ")
	if len(qust) == 0 {
		return "", errors.New("请说出你想问的问题")
	}

	config := openai.DefaultConfig(svcConf.GptKey)
	config.HTTPClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
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
	fmt.Println(err)
	chLen := len(resp.Choices)
	if err != nil || chLen == 0 {
		return "", fmt.Errorf("ChatCompletion len = 0 or error: %v", err)
	}

	return resp.Choices[chLen-1].Message.Content, nil
	// return gpt.qustion(qust, svcConf)
}

func (gpt *Gpt) qustion(qust string, svcConf conf.Config) (string, error) {

	// qust = strings.Trim(strings.ReplaceAll(qust, svcConf.GptKeywords, " "), " ")

	fmt.Printf("开始提问%v", qust)
	cmd := exec.Command("curl", "--max-time", "180", "--request", "POST", `https://api.openai.com/v1/chat/completions`,
		"-H", "Content-Type: application/json",
		"-H", fmt.Sprintf("Authorization: Bearer %v", svcConf.GptKey),
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
