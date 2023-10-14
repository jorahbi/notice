package notice

import (
	"fmt"

	"context"
	"encoding/json"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"

	"github.com/hibiken/asynq"
	"github.com/jorahbi/coco/chain"
	"github.com/jorahbi/notice/internal/aqueue/jobtype"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
)

type WechatNoticeHandler struct {
	svcCtx *svc.ServiceContext
	bot    *openwechat.Bot
	self   *openwechat.Self
}

func NewWechatNoticeHandler(svcCtx *svc.ServiceContext) *WechatNoticeHandler {
	notice := &WechatNoticeHandler{svcCtx: svcCtx}
	notice.start()

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

func (l *WechatNoticeHandler) start() {
	var err error
	fmt.Println("点击确认登录")
	l.bot = openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	// 注册消息处理函数
	l.bot.MessageHandler = l.received
	// 注册登陆二维码回调
	// bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	// if err := bot.Login(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	l.bot.UUIDCallback = l.consoleQrCode
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
		panic(c.Error())
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
	fmt.Println(msg, msg.Content, msg.ToUserName) //filehelper
	payload := client.Payload{
		Fo:   msg.ToUserName,
		Data: msg.Content,
	}
	_, err := l.svcCtx.Client.ReveSend(jobtype.JOB_KEY_GOURMET_RECEIVED_NOTICE, &payload)
	reve := "收到"
	if err != nil {
		reve = err.Error()
	}
	msg.ReplyText(reve)
}

func (l *WechatNoticeHandler) consoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}