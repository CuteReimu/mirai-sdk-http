package miraihttp

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"log/slog"
)

func init() {
	decoder["FriendMessage"] = parseFriendMessage
	decoder["GroupMessage"] = parseGroupMessage
	decoder["TempMessage"] = parseTempMessage
	decoder["StrangerMessage"] = parseStrangerMessage
	decoder["OtherClientMessage"] = parseOtherClientMessage
	decoder["FriendSyncMessage"] = parseFriendSyncMessage
	decoder["GroupSyncMessage"] = parseGroupSyncMessage
	decoder["TempSyncMessage"] = parseTempSyncMessage
	decoder["StrangerSyncMessage"] = parseStrangerSyncMessage
}

// FriendMessage 好友消息
type FriendMessage struct {
	Sender       Friend
	MessageChain []SingleMessage
}

// ListenFriendMessage 监听好友消息
func (b *Bot) ListenFriendMessage(l func(message *FriendMessage) bool) {
	listen(b, "FriendMessage", l)
}

func parseFriendMessage(data gjson.Result) any {
	sender := data.Get("sender")
	if sender.Type != gjson.JSON {
		slog.Error("sender is invalid", "sender", sender)
		return nil
	}
	m := &FriendMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Sender); err != nil {
		slog.Error("json unmarshal failed", "buf", sender.Raw, "error", err)
		return nil
	}
	m.MessageChain = parseMessageChain(data.Get("messageChain").Array())
	return m
}

// GroupMessage 群消息
type GroupMessage struct {
	Sender       Member
	MessageChain []SingleMessage
}

// ListenGroupMessage 监听群消息
func (b *Bot) ListenGroupMessage(l func(message *GroupMessage) bool) {
	listen(b, "GroupMessage", l)
}

func parseGroupMessage(data gjson.Result) any {
	sender := data.Get("sender")
	if sender.Type != gjson.JSON {
		slog.Error("sender is invalid", "sender", sender)
		return nil
	}
	m := &GroupMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Sender); err != nil {
		slog.Error("json unmarshal failed", "buf", sender.Raw, "error", err)
		return nil
	}
	m.MessageChain = parseMessageChain(data.Get("messageChain").Array())
	return m
}

// TempMessage 群临时消息
type TempMessage struct {
	Sender       Member
	MessageChain []SingleMessage
}

// ListenTempMessage 监听群临时消息
func (b *Bot) ListenTempMessage(l func(message *TempMessage) bool) {
	listen(b, "TempMessage", l)
}

func parseTempMessage(data gjson.Result) any {
	sender := data.Get("sender")
	if sender.Type != gjson.JSON {
		slog.Error("sender is invalid", "sender", sender)
		return nil
	}
	m := &TempMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Sender); err != nil {
		slog.Error("json unmarshal failed", "buf", sender.Raw, "error", err)
		return nil
	}
	m.MessageChain = parseMessageChain(data.Get("messageChain").Array())
	return m
}

// StrangerMessage 陌生人消息
type StrangerMessage struct {
	Sender       Friend
	MessageChain []SingleMessage
}

// ListenStrangerMessage 监听陌生人消息
func (b *Bot) ListenStrangerMessage(l func(message *StrangerMessage) bool) {
	listen(b, "StrangerMessage", l)
}

func parseStrangerMessage(data gjson.Result) any {
	sender := data.Get("sender")
	if sender.Type != gjson.JSON {
		slog.Error("sender is invalid", "sender", sender)
		return nil
	}
	m := &StrangerMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Sender); err != nil {
		slog.Error("json unmarshal failed", "buf", sender.Raw, "error", err)
		return nil
	}
	m.MessageChain = parseMessageChain(data.Get("messageChain").Array())
	return m
}

type OtherClient struct {
	Id       int64  `json:"id"`
	Platform string `json:"platform"`
}

// OtherClientMessage 其他客户端消息
type OtherClientMessage struct {
	Sender       OtherClient
	MessageChain []SingleMessage
}

// ListenOtherClientMessage 监听其他客户端消息
func (b *Bot) ListenOtherClientMessage(l func(message *OtherClientMessage) bool) {
	listen(b, "OtherClientMessage", l)
}

func parseOtherClientMessage(data gjson.Result) any {
	sender := data.Get("sender")
	if sender.Type != gjson.JSON {
		slog.Error("sender is invalid", "sender", sender)
		return nil
	}
	m := &OtherClientMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Sender); err != nil {
		slog.Error("json unmarshal failed", "buf", sender.Raw, "error", err)
		return nil
	}
	m.MessageChain = parseMessageChain(data.Get("messageChain").Array())
	return m
}

// FriendSyncMessage 同步好友消息
type FriendSyncMessage struct {
	Subject      Friend
	MessageChain []SingleMessage
}

// ListenFriendSyncMessage 监听同步好友消息
func (b *Bot) ListenFriendSyncMessage(l func(message *FriendSyncMessage) bool) {
	listen(b, "FriendSyncMessage", l)
}

func parseFriendSyncMessage(data gjson.Result) any {
	sender := data.Get("subject")
	if sender.Type != gjson.JSON {
		slog.Error("sender is invalid", "sender", sender)
		return nil
	}
	m := &FriendSyncMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Subject); err != nil {
		slog.Error("json unmarshal failed", "buf", sender.Raw, "error", err)
		return nil
	}
	m.MessageChain = parseMessageChain(data.Get("messageChain").Array())
	return m
}

// GroupSyncMessage 同步群消息
type GroupSyncMessage struct {
	Subject      Group
	MessageChain []SingleMessage
}

// ListenGroupSyncMessage 监听同步群消息
func (b *Bot) ListenGroupSyncMessage(l func(message *GroupSyncMessage) bool) {
	listen(b, "GroupSyncMessage", l)
}

func parseGroupSyncMessage(data gjson.Result) any {
	sender := data.Get("subject")
	if sender.Type != gjson.JSON {
		slog.Error("sender is invalid", "sender", sender)
		return nil
	}
	m := &GroupSyncMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Subject); err != nil {
		slog.Error("json unmarshal failed", "buf", sender.Raw, "error", err)
		return nil
	}
	m.MessageChain = parseMessageChain(data.Get("messageChain").Array())
	return m
}

// TempSyncMessage 同步群临时消息
type TempSyncMessage struct {
	Subject      Member
	MessageChain []SingleMessage
}

// ListenTempSyncMessage 监听同步群临时消息
func (b *Bot) ListenTempSyncMessage(l func(message *TempSyncMessage) bool) {
	listen(b, "TempSyncMessage", l)
}

func parseTempSyncMessage(data gjson.Result) any {
	sender := data.Get("subject")
	if sender.Type != gjson.JSON {
		slog.Error("sender is invalid", "sender", sender)
		return nil
	}
	m := &TempSyncMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Subject); err != nil {
		slog.Error("json unmarshal failed", "buf", sender.Raw, "error", err)
		return nil
	}
	m.MessageChain = parseMessageChain(data.Get("messageChain").Array())
	return m
}

// StrangerSyncMessage 同步好友消息
type StrangerSyncMessage struct {
	Subject      Friend
	MessageChain []SingleMessage
}

// ListenStrangerSyncMessage 监听同步好友消息
func (b *Bot) ListenStrangerSyncMessage(l func(message *StrangerSyncMessage) bool) {
	listen(b, "StrangerSyncMessage", l)
}

func parseStrangerSyncMessage(data gjson.Result) any {
	sender := data.Get("subject")
	if sender.Type != gjson.JSON {
		slog.Error("sender is invalid", "sender", sender)
		return nil
	}
	m := &StrangerSyncMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Subject); err != nil {
		slog.Error("json unmarshal failed", "buf", sender.Raw, "error", err)
		return nil
	}
	m.MessageChain = parseMessageChain(data.Get("messageChain").Array())
	return m
}
