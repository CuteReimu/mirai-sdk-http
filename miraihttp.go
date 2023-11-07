package miraihttp

import (
	"encoding/json"
	"fmt"
	"github.com/CuteReimu/mirai-sdk-http/utils"
	"github.com/gorilla/websocket"
	"sync/atomic"
)

var log = utils.GetModuleLogger("miraihttp")

// WsChannel 连接通道
type WsChannel string

const (
	// WsChannelMessage 推送消息
	WsChannelMessage = "message"

	// WsChannelEvent 推送事件
	WsChannelEvent = "event"

	// WsChannelAll 推送消息及事件
	WsChannelAll = "all"
)

// Connect 连接mirai-api-http
func Connect(host string, port int, channel WsChannel, verifyKey string, qq int64) (*Bot, error) {
	addr := fmt.Sprintf("ws://%s:%d/%s?verifyKey=%s&qq=%d", host, port, channel, verifyKey, qq)
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Errorln("read err:", err)
				return
			}
			log.Debugf("recv: %s\n", message)
		}
	}()
	return &Bot{c: c, QQ: qq}, nil
}

type Bot struct {
	QQ     int64
	c      *websocket.Conn
	syncId atomic.Int64
}

type Request interface {
	GetCommand() string
}

type SubMessage interface {
	Request
	GetSubCommand() string
}

// WriteMessage 发送请求
func (b *Bot) WriteMessage(m Request) {
	msg := &requestMessage{
		SyncId:  b.syncId.Add(1),
		Command: m.GetCommand(),
		Content: m,
	}
	if sub, ok := m.(SubMessage); ok {
		msg.SubCommand = sub.GetSubCommand()
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Errorln("json marshal failed:", err)
		return
	}
	err = b.c.WriteMessage(websocket.TextMessage, buf)
	if err != nil {
		log.Errorln("write err:", err)
		return
	}
	log.Debugf("write: %s\n", string(buf))
}

type requestMessage struct {
	SyncId     int64  `json:"syncId"`
	Command    string `json:"command"`
	SubCommand string `json:"subCommand,omitempty"`
	Content    any    `json:"content,omitempty"`
}
