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

type wechat struct {
	svcCtx *svc.ServiceContext
	bot    *openwechat.Bot
	self   *openwechat.Self
	timer  *time.Timer
}

type WorkerInterface interface {
	Start(ctx context.Context, svcCtx *svc.ServiceContext)
}

var ctx context.Context
var cancel context.CancelFunc

func NewWechat(svcCtx *svc.ServiceContext) *wechat {
	ctx, cancel = context.WithCancel(context.Background())
	wx.svcCtx = svcCtx
	return wx
}

func (l *wechat) Start(works ...WorkerInterface) {
	defer cancel()
	group := threading.NewRoutineGroup()
	group.RunSafe(func() {
		l.waitForSignals(ctx)
	})
	l.start(ctx)
	for _, work := range works {
		work := work
		group.RunSafe(func() {
			work.Start(ctx, l.svcCtx)
		})
	}
	l.bot.Block()
	group.Wait()
}

func (l *wechat) ProcessTask(ctx context.Context, t *asynq.Task) error {
	p := &client.Payload{}
	var err error
	if err = json.Unmarshal(t.Payload(), p); err != nil {
		return err
	}
	l.Send(p)
	return nil
}

func (l *wechat) start(ctx context.Context) {
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
	l.timer.Stop()

}

// waitForSignals waits for signals and handles them.
// It handles SIGTERM, SIGINT, and SIGTSTP.
// SIGTERM and SIGINT will signal the process to exit.
// SIGTSTP will signal the process to stop processing new tasks.
func (l *wechat) waitForSignals(ctx context.Context) {
	fmt.Println("Send signal TERM or INT or TSTP to stop processing new tasks")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	defer l.timer.Stop()

	for {
		select {
		case <-sigs:
			fmt.Println("sigs exit")
			l.bot.Exit()
			return
		case <-ctx.Done():
			fmt.Println("done exit")
			return
		case <-l.timer.C:
			fmt.Printf("timeout login [%v] second", 300)
			cancel()
			os.Exit(0)
		}

	}
}

func (l *wechat) Send(payload *client.Payload) {
	self, err := l.bot.GetCurrentUser()
	fmt.Println(err)
	friends, err := self.Friends()
	friends.SearchByRemarkName(1, payload.Fo).SendText(payload.String())
	fmt.Println(payload)
}

func (l *wechat) received(msg *openwechat.Message) {
	// fmt.Print("消息发送者：")
	// u, e := msg.Sender()
	// if e == nil {
	// 	fmt.Println(e, u.NickName, u.UserName, u.RemarkName)
	// }
	// fmt.Print("消息接收者：")
	// u, e = msg.Receiver()
	// if e == nil {
	// 	fmt.Println(e, u.NickName, u.UserName, u.RemarkName)
	// }
	// fmt.Print("群组消息发送者：")
	// u, e = msg.SenderInGroup()
	// if e == nil {
	// 	fmt.Println(e, u.NickName, u.UserName, u.RemarkName)
	// }
	// fmt.Println(msg.IsSendByFriend(), msg.Owner().NickName)
	//filehelper
	ctx, cancel := context.WithTimeout(context.Background(), 310*time.Second)
	defer cancel()
	for _, event := range event.Events {
		msg.Content = strings.Trim(msg.Content, " ")
		content := []rune(msg.Content)
		if len(content) == 0 {
			return
		}
		recv, err := event.Event(ctx, l.svcCtx, msg)
		l.event(recv, msg, err)

	}
}

func (l *wechat) event(recv string, msg *openwechat.Message, err error) {
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

func (l *wechat) consoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

func (l *wechat) logout(bot *openwechat.Bot) {
	fmt.Println("wechat logout")
	cancel()
}
