package main

import (
	. "github.com/CuteReimu/mirai-sdk-http"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	b, _ := Connect("localhost", 8080, WsChannelAll, "ABCDEFGHIJK", 123456789)
	b.ListenGroupMessage(func(_ string, message *GroupMessage) bool {
		b.WriteMessage(&SendGroupMessage{
			Target:       message.Sender.Group.Id,
			MessageChain: append(MessageChain(&Plain{Text: "你说了：\n"}), message.MessageChain[1:]...),
		})
		return true
	})
	select {}
}
