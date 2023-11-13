package miraihttp

import (
	"encoding/json"
	"github.com/tidwall/gjson"
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

// Friend 好友
type Friend struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
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

func parseFriendMessage(message []byte) any {
	sender := gjson.GetBytes(message, "sender")
	if sender.Type != gjson.JSON {
		log.Errorln("sender is invalid: ", sender)
		return nil
	}
	m := &FriendMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Sender); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	m.MessageChain = parseMessageChain(gjson.GetBytes(message, "messageChain").Array())
	return m
}

// Group 群
type Group struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Permission Perm   `json:"permission"`
}

// Member 群成员
type Member struct {
	Id                 int64  `json:"id"`
	MemberName         string `json:"memberName"`
	SpecialTitle       string `json:"specialTitle"`
	Permission         Perm   `json:"permission"`
	JoinTimestamp      int64  `json:"joinTimestamp"`
	LastSpeakTimestamp int64  `json:"lastSpeakTimestamp"`
	MuteTimeRemaining  int64  `json:"muteTimeRemaining"`
	Group              Group  `json:"group"`
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

func parseGroupMessage(message []byte) any {
	sender := gjson.GetBytes(message, "sender")
	if sender.Type != gjson.JSON {
		log.Errorln("sender is invalid: ", sender)
		return nil
	}
	m := &GroupMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Sender); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	m.MessageChain = parseMessageChain(gjson.GetBytes(message, "messageChain").Array())
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

func parseTempMessage(message []byte) any {
	sender := gjson.GetBytes(message, "sender")
	if sender.Type != gjson.JSON {
		log.Errorln("sender is invalid: ", sender)
		return nil
	}
	m := &TempMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Sender); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	m.MessageChain = parseMessageChain(gjson.GetBytes(message, "messageChain").Array())
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

func parseStrangerMessage(message []byte) any {
	sender := gjson.GetBytes(message, "sender")
	if sender.Type != gjson.JSON {
		log.Errorln("sender is invalid: ", sender)
		return nil
	}
	m := &StrangerMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Sender); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	m.MessageChain = parseMessageChain(gjson.GetBytes(message, "messageChain").Array())
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

func parseOtherClientMessage(message []byte) any {
	sender := gjson.GetBytes(message, "sender")
	if sender.Type != gjson.JSON {
		log.Errorln("sender is invalid: ", sender)
		return nil
	}
	m := &OtherClientMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Sender); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	m.MessageChain = parseMessageChain(gjson.GetBytes(message, "messageChain").Array())
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

func parseFriendSyncMessage(message []byte) any {
	sender := gjson.GetBytes(message, "subject")
	if sender.Type != gjson.JSON {
		log.Errorln("sender is invalid: ", sender)
		return nil
	}
	m := &FriendSyncMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Subject); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	m.MessageChain = parseMessageChain(gjson.GetBytes(message, "messageChain").Array())
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

func parseGroupSyncMessage(message []byte) any {
	sender := gjson.GetBytes(message, "subject")
	if sender.Type != gjson.JSON {
		log.Errorln("sender is invalid: ", sender)
		return nil
	}
	m := &GroupSyncMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Subject); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	m.MessageChain = parseMessageChain(gjson.GetBytes(message, "messageChain").Array())
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

func parseTempSyncMessage(message []byte) any {
	sender := gjson.GetBytes(message, "subject")
	if sender.Type != gjson.JSON {
		log.Errorln("sender is invalid: ", sender)
		return nil
	}
	m := &TempSyncMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Subject); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	m.MessageChain = parseMessageChain(gjson.GetBytes(message, "messageChain").Array())
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

func parseStrangerSyncMessage(message []byte) any {
	sender := gjson.GetBytes(message, "subject")
	if sender.Type != gjson.JSON {
		log.Errorln("sender is invalid: ", sender)
		return nil
	}
	m := &StrangerSyncMessage{}
	if err := json.Unmarshal([]byte(sender.Raw), &m.Subject); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	m.MessageChain = parseMessageChain(gjson.GetBytes(message, "messageChain").Array())
	return m
}
