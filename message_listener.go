package miraihttp

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

func init() {
	parser["FriendMessage"] = parseFriendMessage
	parser["GroupMessage"] = parseGroupMessage
}

type Friend struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
}

type FriendMessage struct {
	Sender       Friend
	MessageChain []SingleMessage
}

func (b *Bot) ListenFriendMessage(l func(*FriendMessage) bool) {
	listenMessage(b, "FriendMessage", l)
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

type Group struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Permission Perm   `json:"permission"`
}

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

type GroupMessage struct {
	Sender       Member
	MessageChain []SingleMessage
}

func (b *Bot) ListenGroupMessage(l func(*GroupMessage) bool) {
	listenMessage(b, "GroupMessage", l)
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
