package miraihttp

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

func init() {
	decoder["NewFriendRequestEvent"] = parseNewFriendRequestEvent
	decoder["MemberJoinRequestEvent"] = parseMemberJoinRequestEvent
	decoder["BotInvitedJoinGroupRequestEvent"] = parseBotInvitedJoinGroupRequestEvent
}

// NewFriendRequestEvent 添加好友申请
type NewFriendRequestEvent struct {
	EventId int64  `json:"eventId"` // 事件标识，响应该事件时的标识
	QQ      int64  `json:"fromId"`  // 申请人QQ号
	Group   int64  `json:"groupId"` // 申请人如果通过某个群添加好友，该项为该群群号；否则为0
	Nick    string `json:"nick"`    // 申请人的昵称或群名片
	Message string `json:"message"` // 申请消息
}

// ListenNewFriendRequestEvent 监听添加好友申请
func (b *Bot) ListenNewFriendRequestEvent(l func(message *NewFriendRequestEvent) bool) {
	listen(b, "NewFriendRequestEvent", l)
}

func parseNewFriendRequestEvent(data gjson.Result) any {
	if data.Type != gjson.JSON {
		log.Errorln("data is invalid: ", data)
		return nil
	}
	m := &NewFriendRequestEvent{}
	if err := json.Unmarshal([]byte(data.Raw), m); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	return m
}

// MemberJoinRequestEvent 用户入群申请（Bot需要有管理员权限）
type MemberJoinRequestEvent struct {
	EventId   int64  `json:"eventId"`   // 事件标识，响应该事件时的标识
	QQ        int64  `json:"fromId"`    // 申请人QQ号
	Group     int64  `json:"groupId"`   // 申请人申请入群的群号
	GroupName string `json:"groupName"` // 申请人申请入群的群名称
	Nick      string `json:"nick"`      // 申请人的昵称或群名片
	Message   string `json:"message"`   // 申请消息
	InvitorId int64  `json:"invitorId"` // 邀请人，可能没有
}

// ListenMemberJoinRequestEvent 监听用户入群申请（Bot需要有管理员权限）
func (b *Bot) ListenMemberJoinRequestEvent(l func(message *MemberJoinRequestEvent) bool) {
	listen(b, "MemberJoinRequestEvent", l)
}

func parseMemberJoinRequestEvent(data gjson.Result) any {
	if data.Type != gjson.JSON {
		log.Errorln("data is invalid: ", data)
		return nil
	}
	m := &MemberJoinRequestEvent{}
	if err := json.Unmarshal([]byte(data.Raw), m); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	return m
}

// BotInvitedJoinGroupRequestEvent Bot被邀请入群申请
type BotInvitedJoinGroupRequestEvent struct {
	EventId   int64  `json:"eventId"`   // 事件标识，响应该事件时的标识
	QQ        int64  `json:"fromId"`    // 邀请人QQ号
	Group     int64  `json:"groupId"`   // 被邀请进入群的群号
	GroupName string `json:"groupName"` // 被邀请进入群的群名称
	Nick      string `json:"nick"`      // 邀请人（好友）的昵称
	Message   string `json:"message"`   // 邀请消息
}

// ListenBotInvitedJoinGroupRequestEvent 监听Bot被邀请入群申请
func (b *Bot) ListenBotInvitedJoinGroupRequestEvent(l func(message *BotInvitedJoinGroupRequestEvent) bool) {
	listen(b, "BotInvitedJoinGroupRequestEvent", l)
}

func parseBotInvitedJoinGroupRequestEvent(data gjson.Result) any {
	if data.Type != gjson.JSON {
		log.Errorln("data is invalid: ", data)
		return nil
	}
	m := &BotInvitedJoinGroupRequestEvent{}
	if err := json.Unmarshal([]byte(data.Raw), m); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	return m
}
