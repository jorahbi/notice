package notice

import (
	"fmt"
	"strings"
	"time"

	// "net/http"

	"context"
	"encoding/json"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"

	"github.com/hibiken/asynq"
	"github.com/jorahbi/coco/chain"

	// "github.com/jorahbi/notice/internal/received"
	// "github.com/jorahbi/notice/internal/event"
	"github.com/jorahbi/notice/internal/event"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
	"github.com/samber/lo"
)

type WechatNoticeHandler struct {
	cancal func()
	svcCtx *svc.ServiceContext
	bot    *openwechat.Bot
	self   *openwechat.Self
}

func NewWechatNoticeHandler(ctx context.Context, svcCtx *svc.ServiceContext, cancal func()) *WechatNoticeHandler {
	notice := &WechatNoticeHandler{svcCtx: svcCtx, cancal: cancal}
	notice.start(ctx)

	return notice
}

func (l *WechatNoticeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	p := &client.Payload{}
	var err error
	if err = json.Unmarshal(t.Payload(), p); err != nil {
		return err
	}
	l.send(p)
	return nil
}

func (l *WechatNoticeHandler) start(ctx context.Context) {
	var err error
	fmt.Println("点击确认登录")
	l.bot = openwechat.DefaultBot(openwechat.Desktop, openwechat.WithContextOption(ctx)) // 桌面模式
	// 注册消息处理函数
	l.bot.MessageHandler = l.received
	// 注册登陆二维码回调
	// bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	// if err := bot.Login(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	l.bot.UUIDCallback = l.consoleQrCode
	l.bot.LogoutCallBack = l.logout
	reloadStorage := openwechat.NewFileHotReloadStorage("etc/storage.json")
	defer reloadStorage.Close()
	c := chain.NewChain()
	c.Apply(func() error { // 登陆
		return l.bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())
	})
	c.Apply(func() error {
		l.self, err = l.bot.GetCurrentUser()
		return err
	})
	err = c.Error()
	if err != nil {
		panic(err)
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	// bot.Block()
}

func (l *WechatNoticeHandler) send(payload *client.Payload) {
	self, err := l.bot.GetCurrentUser()
	fmt.Println(err)
	friends, err := self.Friends()
	friends.SearchByRemarkName(1, payload.Fo).SendText(payload.String())
	fmt.Println(payload)
}

func (l *WechatNoticeHandler) received(msg *openwechat.Message) {
	//filehelper
	ctx, cancel := context.WithTimeout(context.Background(), 310*time.Second)
	defer cancel()
	recv, err := event.MustEventFactory(event.EVENT_KEY_GPT).Event(ctx, l.svcCtx, client.Payload{Data: msg.Content})
	if err != nil {
		recv = fmt.Sprintf("%v%v", recv, err.Error())
	}
	recv = strings.Trim(recv, "")
	content := []rune(strings.Trim(recv, ""))
	if len(content) == 0 {
		return
	}
	for _, val := range lo.Chunk[rune](content, 500) {
		msg.ReplyText(string(val))
	}
}

func (l *WechatNoticeHandler) consoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

func (l *WechatNoticeHandler) logout(bot *openwechat.Bot) {
	fmt.Println("logout")
	l.cancal()
}
