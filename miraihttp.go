package miraihttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CuteReimu/goutil"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"log/slog"
	"runtime/debug"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

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
	log := slog.With("addr", addr)
	log.Info("Dialing")
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Error("Connect failed")
		return nil, err
	}
	log.Info("Connected successfully")
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
				log.Error("read error", "error", err)
				return
			}
			log.Debug("recv: " + string(message))
			if !gjson.ValidBytes(message) {
				log.Error("invalid json message: " + string(message))
				continue
			}
			syncId := gjson.GetBytes(message, "syncId").String()
			data := gjson.GetBytes(message, "data")
			if data.Type != gjson.JSON {
				log.Error("invalid json message: " + string(message))
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
			if h, ok := b.handler.Load(messageType); ok {
				if p := decoder[messageType]; p == nil {
					log.Error("cannot find message decoder: " + messageType)
				} else if m := p(data); m != nil {
					fun := func() {
						defer func() {
							if r := recover(); r != nil {
								log.Error("panic recovered", "error", r, "stack", string(debug.Stack()))
							}
						}()
						h.(*handler).handle(m)
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

type handler struct {
	fs []listenHandler
}

func (h *handler) Append(f listenHandler) *handler {
	if h == nil {
		return &handler{fs: []listenHandler{f}}
	}
	return &handler{fs: append(slices.Clip(h.fs), []listenHandler{f}...)}
}

func (h *handler) handle(m any) {
	if h == nil {
		return
	}
	for _, f := range h.fs {
		if !f(m) {
			break
		}
	}
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
	log := slog.With("command", command, "subCommand", subCommand)
	syncId := strconv.FormatInt(msg.SyncId, 10)
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Error("json marshal failed", "error", err)
		return gjson.Result{}, err
	}
	ch := make(chan gjson.Result, 1)
	b.syncIdMap.Store(syncId, ch)
	err = b.c.WriteMessage(websocket.TextMessage, buf)
	if err != nil {
		log.Error("write error", "error", err)
		return gjson.Result{}, err
	}
	log.Debug("write: " + string(buf))
	timeoutTimer := time.AfterFunc(5*time.Second, func() {
		if ch, ok := b.syncIdMap.LoadAndDelete(syncId); ok {
			close(ch.(chan gjson.Result))
		}
	})
	result, ok := <-ch
	if !ok {
		log.Error("request timeout")
		return gjson.Result{}, errors.New("request timeout")
	}
	timeoutTimer.Stop()
	code := result.Get("code").Int()
	if code != 0 {
		e := fmt.Sprint("Non-zero code: ", code, ", error message: ", result.Get("msg"))
		log.Error(e)
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
	var fs *handler
	if f, ok := b.handler.Load(key); ok {
		fs = f.(*handler)
	}
	fs2 := fs.Append(func(m any) bool { return l(m.(M)) })
	if !b.handler.CompareAndSwap(key, fs, fs2) {
		panic("don't call listen concurrently")
	}
}
