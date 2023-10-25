package notice

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
	"golang.org/x/sys/unix"

	"github.com/hibiken/asynq"
	"github.com/jorahbi/notice/internal/event"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/threading"
)

type WechatNoticeHandler struct {
	svcCtx *svc.ServiceContext
	bot    *openwechat.Bot
	self   *openwechat.Self
}

type WorkerInterface interface {
	Start(ctx context.Context, svcCtx *svc.ServiceContext)
}

var ctx context.Context
var cancel context.CancelFunc
var Wechat asynq.Handler

func Start(svcCtx *svc.ServiceContext, works ...WorkerInterface) {
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	chat := &WechatNoticeHandler{svcCtx: svcCtx}
	works = append(works, chat)
	group := threading.NewRoutineGroup()
	chat.bot = openwechat.DefaultBot(openwechat.Desktop, openwechat.WithContextOption(ctx)) // 桌面模式
	Wechat = chat
	for _, work := range works {
		work := work
		group.RunSafe(func() {
			work.Start(ctx, svcCtx)
		})
	}
	group.RunSafe(func() {
		chat.waitForSignals(ctx)
	})

	group.Wait()
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

func (l WechatNoticeHandler) Start(ctx context.Context, svcCtx *svc.ServiceContext) {
	var err error
	fmt.Println("点击确认登录")

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
	err = l.bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())
	if err != nil {
		fmt.Println("登录失败")
		return
	}
	l.self, err = l.bot.GetCurrentUser()
	if err != nil {
		fmt.Println("获取当前登录")
		return
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	l.bot.Block()
}

// waitForSignals waits for signals and handles them.
// It handles SIGTERM, SIGINT, and SIGTSTP.
// SIGTERM and SIGINT will signal the process to exit.
// SIGTSTP will signal the process to stop processing new tasks.
func (l *WechatNoticeHandler) waitForSignals(ctx context.Context) {
	fmt.Println("Send signal TERM or INT or TSTP to stop processing new tasks")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	var timeout time.Duration = 300
	timer := time.NewTimer(timeout * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-sigs:
			l.bot.Exit()
			return
		case <-ctx.Done():
			return
		case <-timer.C:
			fmt.Printf("timeout login [%v] second", timeout)
			cancel()
			return
		}

	}
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
	for _, event := range event.Events {
		msg.Content = strings.Trim(msg.Content, " ")
		content := []rune(msg.Content)
		if len(content) == 0 {
			return
		}
		recv, err := event.Event(ctx, l.svcCtx, &client.Payload{Data: msg.Content})
		l.event(recv, msg, err)

	}
}

func (l *WechatNoticeHandler) event(recv string, msg *openwechat.Message, err error) {
	if err != nil {
		recv = fmt.Sprintf("%v%v", recv, err.Error())
	}
	content := []rune(strings.Trim(recv, " "))
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
	fmt.Println("wechat logout")
	cancel()
}
