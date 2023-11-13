package main

import (
	. "github.com/CuteReimu/mirai-sdk-http"
	"github.com/CuteReimu/mirai-sdk-http/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	utils.InitLogger("./logs", utils.LogDebugLevel, utils.LogWithStack)
	log := utils.GetModuleLogger("robot")
	b, _ := Connect("localhost", 8080, WsChannelAll, "ABCDEFGHIJK", 123456789, false)
	b.ListenGroupMessage(func(message *GroupMessage) bool {
		_, err := b.SendGroupMessage(message.Sender.Group.Id, 0,
			append(MessageChain(&Plain{Text: "你说了：\n"}), message.MessageChain[1:]...),
		)
		if err != nil {
			log.Println("发送失败: ", err)
		}
		return true
	})
	select {}
}
