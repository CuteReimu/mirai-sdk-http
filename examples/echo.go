package main

import (
	. "github.com/CuteReimu/mirai-sdk-http"
	"github.com/CuteReimu/mirai-sdk-http/utils"
	"log/slog"
)

func main() {
	utils.InitLogger("./logs", slog.LevelDebug)
	b, _ := Connect("localhost", 8080, WsChannelAll, "ABCDEFGHIJK", 123456789, false)
	b.ListenGroupMessage(func(message *GroupMessage) bool {
		_, err := b.SendGroupMessage(message.Sender.Group.Id, 0,
			append(MessageChain(&Plain{Text: "你说了：\n"}), message.MessageChain[1:]...),
		)
		if err != nil {
			slog.Error("发送失败", "error", err)
		}
		return true
	})
	select {}
}
