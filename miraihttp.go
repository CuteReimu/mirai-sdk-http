package miraihttp

import (
	"encoding/json"
	"fmt"
	"github.com/CuteReimu/mirai-sdk-http/utils"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"sync"
	"sync/atomic"
)

var log = utils.GetModuleLogger("miraihttp")

// WsChannel 连接通道
type WsChannel string

const (
	// WsChannelMessage 推送消息
	WsChannelMessage = "model"

	// WsChannelEvent 推送事件
	WsChannelEvent = "event"

	// WsChannelAll 推送消息及事件
	WsChannelAll = "all"
)

// Connect 连接mirai-api-http
func Connect(host string, port int, channel WsChannel, verifyKey string, qq int64) (*Bot, error) {
	addr := fmt.Sprintf("ws://%s:%d/%s?verifyKey=%s&qq=%d", host, port, channel, verifyKey, qq)
	log.Infoln("Dialing", addr)
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}
	log.Infoln("Connected successfully")
	b := &Bot{c: c, QQ: qq}
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Errorln("read err:", err)
				return
			}
			log.Debugf("recv: %s\n", message)
			if !gjson.ValidBytes(message) {
				log.Errorln("invalid json message: ", string(message))
				continue
			}
			syncId := gjson.GetBytes(message, "syncId").String()
			data := gjson.GetBytes(message, "data")
			if data.Type != gjson.JSON {
				log.Errorln("invalid json message: ", string(message))
				continue
			}
			messageType := data.Get("type").String()
			if f, ok := b.handler.Load(messageType); ok {
				if p := parser[messageType]; p == nil {
					log.Errorln("cannot find message parser:", messageType)
				} else if m := p([]byte(data.Raw)); m != nil {
					for _, handler := range f.([]ListenMessageHandler) {
						if !handler(syncId, m) {
							break
						}
					}
				}
			}
		}
	}()
	return b, nil
}

type Bot struct {
	QQ      int64
	c       *websocket.Conn
	syncId  atomic.Int64
	handler sync.Map
}

type Request interface {
	GetCommand() string
}

type SubMessage interface {
	Request
	GetSubCommand() string
}

// WriteMessage 发送请求
func (b *Bot) WriteMessage(m Request) (syncId int64) {
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
		return -1
	}
	err = b.c.WriteMessage(websocket.TextMessage, buf)
	if err != nil {
		log.Errorln("write err:", err)
		return -1
	}
	log.Debugf("write: %s\n", string(buf))
	return msg.SyncId
}

type requestMessage struct {
	SyncId     int64  `json:"syncId"`
	Command    string `json:"command"`
	SubCommand string `json:"subCommand,omitempty"`
	Content    any    `json:"content,omitempty"`
}

var parser = make(map[string]func(message []byte) any)

type ListenMessageHandler func(syncId string, message any) bool

func listenMessage[M any](b *Bot, key string, l func(syncId string, message M) bool) {
	var fs []ListenMessageHandler
	if f, ok := b.handler.Load(key); ok {
		fs = f.([]ListenMessageHandler)
	}
	fs = append(fs, func(syncId string, m any) bool {
		return l(syncId, m.(M))
	})
	b.handler.Store(key, fs)
}
