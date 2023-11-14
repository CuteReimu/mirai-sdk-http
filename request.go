package miraihttp

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// About 获取插件版本号
func (b *Bot) About() (string, error) {
	result, err := b.request("sendGroupMessage", "", nil)
	if err != nil {
		return "", err
	}
	return result.Get("data.version").String(), nil
}

// BotList 获取登录账号
func (b *Bot) BotList() ([]int64, error) {
	result, err := b.request("botList", "", nil)
	if err != nil {
		return nil, err
	}
	data := result.Get("data").Array()
	bots := make([]int64, 0, len(data))
	for _, r := range data {
		bots = append(bots, r.Int())
	}
	return bots, nil
}

// MessageFromId 通过messageId获取消息，target-好友或QQ群，视情况返回 FriendMessage, GroupMessage, TempMessage, StrangerMessage
func (b *Bot) MessageFromId(messageId, target int64) (any, error) {
	result, err := b.request("messageFromId", "", &struct {
		MessageId int64 `json:"messageId"`
		Target    int64 `json:"target"`
	}{messageId, target})
	if err != nil {
		return nil, err
	}
	data := result.Get("data")
	if data.Type != gjson.JSON {
		e := fmt.Sprint("invalid json message: ", result)
		log.Errorln(e)
		return nil, errors.New(e)
	}
	messageType := data.Get("type").String()
	if p := decoder[messageType]; p != nil {
		if m := p([]byte(data.Raw)); m != nil {
			return m, nil
		}
	}
	e := fmt.Sprint("decode message failed:", data.Raw)
	log.Errorln(e)
	return nil, errors.New(e)
}

// SendFriendMessage 发送好友消息，qq-目标好友的QQ号，quote-引用回复的消息，messageChain-发送的内容，返回消息id
func (b *Bot) SendFriendMessage(qq, quote int64, messageChain []SingleMessage) (int64, error) {
	result, err := b.request("sendFriendMessage", "", &struct {
		Target       int64           `json:"target"`
		Quote        int64           `json:"quote,omitempty"`
		MessageChain []SingleMessage `json:"messageChain"`
	}{qq, quote, messageChain})
	if err != nil {
		return 0, err
	}
	return result.Get("messageId").Int(), nil
}

// SendGroupMessage 发送群消息，group-群号，quote-引用回复的消息，messageChain-发送的内容，返回消息id
func (b *Bot) SendGroupMessage(group, quote int64, messageChain []SingleMessage) (int64, error) {
	result, err := b.request("sendGroupMessage", "", &struct {
		Target       int64           `json:"target"`
		Quote        int64           `json:"quote,omitempty"`
		MessageChain []SingleMessage `json:"messageChain"`
	}{group, quote, messageChain})
	if err != nil {
		return 0, err
	}
	return result.Get("messageId").Int(), nil
}

// SendTempMessage 发送临时会话消息，qq-临时会话对象QQ号，group-临时会话群号，quote-引用回复的消息，messageChain-发送的内容，返回消息id
func (b *Bot) SendTempMessage(qq, group, quote int64, messageChain []SingleMessage) (int64, error) {
	result, err := b.request("sendTempMessage", "", &struct {
		QQ           int64           `json:"qq"`
		Group        int64           `json:"group"`
		Quote        int64           `json:"quote,omitempty"`
		MessageChain []SingleMessage `json:"messageChain"`
	}{qq, group, quote, messageChain})
	if err != nil {
		return 0, err
	}
	return result.Get("messageId").Int(), nil
}

// SendNudge 发送头像戳一戳消息，qq-戳谁，subject-这条消息发到哪（好友/群），kind-上下文类型
func (b *Bot) SendNudge(qq, subject int64, kind Kind) error {
	_, err := b.request("sendNudge", "", &struct {
		Target  int64 `json:"target"`
		Subject int64 `json:"subject"`
		Kind    Kind  `json:"kind"`
	}{qq, subject, kind})
	return err
}

// Recall 撤回消息，target-撤回哪的消息（好友/群），messageId-需要撤回的消息的messageId
func (b *Bot) Recall(target, messageId int64) error {
	_, err := b.request("recall", "", &struct {
		Target    int64 `json:"target"`
		MessageId int64 `json:"messageId"`
	}{target, messageId})
	return err
}

// RoamingMessages 获取漫游消息，timeStart和timeEnd为开始和结束的时间戳，单位为秒。qq为查询的对象QQ，目前仅支持好友漫游消息。
//
// 返回数组的元素为 FriendMessage, GroupMessage, TempMessage, StrangerMessage
func (b *Bot) RoamingMessages(timeStart, timeEnd, qq int64) ([]any, error) {
	result, err := b.request("roamingMessages", "", &struct {
		TimeStart int64 `json:"timeStart"`
		TimeEnd   int64 `json:"timeEnd"`
		Target    int64 `json:"target"`
	}{timeStart, timeEnd, qq})
	if err != nil {
		return nil, err
	}
	dataArray := result.Get("data").Array()
	retArray := make([]any, 0, len(dataArray))
	for _, data := range dataArray {
		if data.Type != gjson.JSON {
			e := fmt.Sprint("invalid json message: ", result)
			log.Errorln(e)
			return nil, errors.New(e)
		}
		messageType := data.Get("type").String()
		if p := decoder[messageType]; p != nil {
			if m := p([]byte(data.Raw)); m != nil {
				retArray = append(retArray, m)
				continue
			}
		}
		e := fmt.Sprint("decode message failed:", data.Raw)
		log.Errorln(e)
		return nil, errors.New(e)
	}
	return retArray, nil
}

// Mute 禁言群成员（需要有相关限权），group-群，qq-被禁言的人，time-时间，单位秒，最多30天
func (b *Bot) Mute(group, qq, time int64) error {
	_, err := b.request("mute", "", &struct {
		Target   int64 `json:"target"`
		MemberId int64 `json:"memberId"`
		Time     int64 `json:"time"`
	}{group, qq, time})
	return err
}

// Unmute 解除禁言群成员（需要有相关限权），group-群，qq-解除禁言的人
func (b *Bot) Unmute(group, qq int64) error {
	_, err := b.request("unmute", "", &struct {
		Target   int64 `json:"target"`
		MemberId int64 `json:"memberId"`
	}{group, qq})
	return err
}

// Kick 移除群成员（需要有相关限权），group-群，qq-移除的人，block-移除后是否拉黑，msg-信息
func (b *Bot) Kick(group, qq int64, block bool, msg string) error {
	_, err := b.request("kick", "", &struct {
		Target   int64  `json:"target"`
		MemberId int64  `json:"memberId"`
		Block    bool   `json:"block"`
		Msg      string `json:"msg"`
	}{group, qq, block, msg})
	return err
}

// Quit 退出群聊（自己不能是群主）
func (b *Bot) Quit(group int64) error {
	_, err := b.request("quit", "", &struct {
		Target int64 `json:"target"`
	}{group})
	return err
}

// MuteAll 全体禁言（需要有相关限权）
func (b *Bot) MuteAll(group int64) error {
	_, err := b.request("muteAll", "", &struct {
		Target int64 `json:"target"`
	}{group})
	return err
}

// UnmuteAll 解除全体禁言（需要有相关限权）
func (b *Bot) UnmuteAll(group int64) error {
	_, err := b.request("unmuteAll", "", &struct {
		Target int64 `json:"target"`
	}{group})
	return err
}

// SetEssence 设置群精华消息（需要有相关限权）
func (b *Bot) SetEssence(group, messageId int64) error {
	_, err := b.request("setEssence", "", &struct {
		Target    int64 `json:"target"`
		MessageId int64 `json:"messageId"`
	}{group, messageId})
	return err
}

// GroupConfig 群设置
type GroupConfig struct {
	Name              string `json:"name"`
	Announcement      string `json:"announcement"`
	ConfessTalk       bool   `json:"confessTalk"`
	AllowMemberInvite bool   `json:"allowMemberInvite"`
	AutoApprove       bool   `json:"autoApprove"`
	AnonymousChat     bool   `json:"anonymousChat"`

	// MuteAll 是否禁言，修改群设置时不要填这个字段，而应该用 Bot.MuteAll(group) 方法
	MuteAll bool `json:"muteAll,omitempty"`
}

// GetGroupConfig 获取群设置
func (b *Bot) GetGroupConfig(group int64) (*GroupConfig, error) {
	result, err := b.request("groupConfig", "get", &struct {
		Target int64 `json:"target"`
	}{group})
	if err != nil {
		return nil, err
	}
	groupConfig := &GroupConfig{}
	if err = json.Unmarshal([]byte(result.Raw), groupConfig); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		log.Errorln(e)
		return nil, err
	}
	return groupConfig, nil
}

// UpdateGroupConfig 修改群设置（需要有相关限权）
func (b *Bot) UpdateGroupConfig(group int64, groupConfig *GroupConfig) error {
	_, err := b.request("groupConfig", "update", &struct {
		Target int64        `json:"target"`
		Config *GroupConfig `json:"config"`
	}{group, groupConfig})
	return err
}

// GetMemberInfo 获取群员设置
func (b *Bot) GetMemberInfo(group, qq int64) (*Member, error) {
	result, err := b.request("memberInfo", "get", &struct {
		Target   int64 `json:"target"`
		MemberId int64 `json:"memberId"`
	}{group, qq})
	if err != nil {
		return nil, err
	}
	member := &Member{}
	if err = json.Unmarshal([]byte(result.Raw), member); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		log.Errorln(e)
		return nil, err
	}
	return member, nil
}

// UpdateMemberInfo 修改群员设置（需要有相关限权），name-群昵称，specialTitle-群头衔，这两项都是选填
func (b *Bot) UpdateMemberInfo(group, qq int64, name, specialTitle string) error {
	type Info struct {
		Name         string `json:"name,omitempty"`
		SpecialTitle string `json:"specialTitle,omitempty"`
	}
	_, err := b.request("memberInfo", "update", &struct {
		Target   int64 `json:"target"`
		MemberId int64 `json:"memberId"`
		Info     Info  `json:"info"`
	}{group, qq, Info{name, specialTitle}})
	return err
}

// MemberAdmin 修改群员管理员（需要有群主限权），assign-是否设置为管理员
func (b *Bot) MemberAdmin(group, qq int64, assign bool) error {
	_, err := b.request("memberAdmin", "", &struct {
		Target   int64 `json:"target"`
		MemberId int64 `json:"memberId"`
		Assign   bool  `json:"assign"`
	}{group, qq, assign})
	return err
}
