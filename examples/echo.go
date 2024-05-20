package main

import (
	. "github.com/CuteReimu/mirai-sdk-http"
	"golang.org/x/time/rate"
	"log/slog"
)

func main() {
	b, err := Connect("localhost", 8080, WsChannelAll, "ABCDEFGHIJK", 123456789, false)
	if err != nil {
		panic(err)
	}
	// 设置限流策略为：令牌桶容量为10，每秒放入一个令牌，超过的消息直接丢弃
	b.SetLimiter("drop", rate.NewLimiter(1, 10))
	b.ListenGroupMessage(func(message *GroupMessage) bool {
		var ret MessageChain
		ret = append(ret, &Plain{Text: "你说了：\n"})
		ret = append(ret, message.MessageChain[1:]...) // 第一个元素原先一定是 Source ，直接排除掉即可
		_, err := b.SendGroupMessage(message.Sender.Group.Id, 0, ret)
		if err != nil {
			slog.Error("发送失败", "error", err)
		}
		return true
	})
	select {}
}
