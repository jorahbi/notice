package notice

import (
	"fmt"

	"context"
	"encoding/json"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"

	"github.com/hibiken/asynq"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
)

type WechatNoticeHandler struct {
	svcCtx *svc.ServiceContext
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
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		fmt.Println(msg, msg.Content, msg.ToUserName) //filehelper
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	// 注册登陆二维码回调
	// bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	// if err := bot.Login(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	bot.UUIDCallback = l.consoleQrCode
	reloadStorage := openwechat.NewFileHotReloadStorage("etc/storage.json")
	defer reloadStorage.Close()
	err := bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())

	// 登陆
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)
	// fmt.Println(friends.GetByUsername("filehelper").SendText("filehelper"))
	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	// bot.Block()
}

func (l *WechatNoticeHandler) send(payload *client.Payload) {
	fmt.Println(payload)
}

func (l *WechatNoticeHandler) consoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}
