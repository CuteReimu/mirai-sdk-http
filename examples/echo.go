package main

import (
	. "github.com/CuteReimu/mirai-sdk-http"
)

func main() {
	b, _ := Connect("localhost", 8080, WsChannelAll, "ABCDEFGHIJK", 123456789)
	b.ListenGroupMessage(func(message *GroupMessage) bool {
		b.WriteMessage(&SendGroupMessage{
			Target: message.Sender.Group.Id,
			MessageChain: MessageChain(
				&Plain{Text: "Mirai牛逼"},
			),
		})
		return true
	})
}
