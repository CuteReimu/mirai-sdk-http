package miraihttp

import (
	"encoding/json"
	"fmt"
	"github.com/CuteReimu/goutil"
	"github.com/CuteReimu/mirai-sdk-http/utils"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
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
//
// concurrentEvent 参数如果是true，表示采用并发方式处理事件和消息，由调用者自行解决并发问题。
// 如果是false表示用单线程处理事件和消息，调用者无需关心并发问题。
func Connect(host string, port int, channel WsChannel, verifyKey string, qq int64, concurrentEvent bool) (*Bot, error) {
	addr := fmt.Sprintf("ws://%s:%d/%s?verifyKey=%s&qq=%d", host, port, channel, verifyKey, qq)
	log.Infoln("Dialing", addr)
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Errorln("Connect failed")
		return nil, err
	}
	log.Infoln("Connected successfully")
	b := &Bot{QQ: qq, c: c}
	if !concurrentEvent {
		b.eventChan = goutil.NewBlockingQueue[func()]()
		go func() {
			for {
				b.eventChan.Take()()
			}
		}()
	}
	go func() {
		for {
			t, message, err := c.ReadMessage()
			if t != websocket.TextMessage {
				continue
			}
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
			if len(syncId) > 0 && syncId[0] != '-' {
				if ch, ok := b.syncIdMap.LoadAndDelete(syncId); ok {
					ch0 := ch.(chan gjson.Result)
					ch0 <- data
					close(ch0)
				}
				continue
			}
			messageType := data.Get("type").String()
			if f, ok := b.handler.Load(messageType); ok {
				if p := decoder[messageType]; p == nil {
					log.Errorln("cannot find message decoder:", messageType)
				} else if m := p(data); m != nil {
					fun := func() {
						defer func() {
							if r := recover(); r != nil {
								log.Errorln("panic recovered: ", r)
							}
						}()
						for _, handler := range f.([]listenHandler) {
							if !handler(m) {
								break
							}
						}
					}
					if b.eventChan == nil {
						go fun()
					} else {
						b.eventChan.Put(fun)
					}
				}
			}
		}
	}()
	return b, nil
}

type Bot struct {
	QQ        int64
	c         *websocket.Conn
	syncId    atomic.Int64
	handler   sync.Map
	syncIdMap sync.Map
	eventChan *goutil.BlockingQueue[func()]
}

// request 发送请求
func (b *Bot) request(command, subCommand string, m any) (gjson.Result, error) {
	msg := &requestMessage{
		SyncId:     b.syncId.Add(1),
		Command:    command,
		SubCommand: subCommand,
		Content:    m,
	}
	syncId := strconv.FormatInt(msg.SyncId, 10)
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Errorln("json marshal failed:", err)
		return gjson.Result{}, err
	}
	ch := make(chan gjson.Result, 1)
	b.syncIdMap.Store(syncId, ch)
	err = b.c.WriteMessage(websocket.TextMessage, buf)
	if err != nil {
		log.Errorln("write err:", err)
		return gjson.Result{}, err
	}
	log.Debugf("write: %s\n", string(buf))
	time.AfterFunc(5*time.Second, func() {
		if ch, ok := b.syncIdMap.LoadAndDelete(syncId); ok {
			close(ch.(chan gjson.Result))
		}
	})
	result, ok := <-ch
	if !ok {
		log.Errorln("request timeout")
		return gjson.Result{}, errors.New("request timeout")
	}
	code := result.Get("code").Int()
	if code != 0 {
		e := fmt.Sprint("Non-zero code: ", code, ", error message: ", result.Get("msg"))
		log.Errorln(e)
		return gjson.Result{}, errors.New(e)
	}
	return result, nil
}

type requestMessage struct {
	SyncId     int64  `json:"syncId"`
	Command    string `json:"command"`
	SubCommand string `json:"subCommand,omitempty"`
	Content    any    `json:"content,omitempty"`
}

var decoder = make(map[string]func(data gjson.Result) any)

type listenHandler func(message any) bool

func listen[M any](b *Bot, key string, l func(message M) bool) {
	var fs []listenHandler
	if f, ok := b.handler.Load(key); ok {
		fs = f.([]listenHandler)
	}
	fs = append(fs, func(m any) bool {
		return l(m.(M))
	})
	b.handler.Store(key, fs)
}
