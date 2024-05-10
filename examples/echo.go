package main

import (
	. "github.com/CuteReimu/mirai-sdk-http"
	"log/slog"
)

func main() {
	b, _ := Connect("localhost", 8080, WsChannelAll, "ABCDEFGHIJK", 123456789, false)
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
