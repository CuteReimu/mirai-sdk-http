package miraihttp

import (
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

// SendNudge 发送头像戳一戳消息，target-戳谁，subject-这条消息发到哪（好友/群），kind-上下文类型
func (b *Bot) SendNudge(target, subject int64, kind Kind) error {
	_, err := b.request("sendNudge", "", &struct {
		Target  int64 `json:"target"`
		Subject int64 `json:"subject"`
		Kind    Kind  `json:"kind"`
	}{target, subject, kind})
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

// RoamingMessages 获取漫游消息，timeStart和timeEnd为开始和结束的时间戳，单位为秒。target为查询的对象，目前仅支持好友漫游消息。
//
// 返回数组的元素为 FriendMessage, GroupMessage, TempMessage, StrangerMessage
func (b *Bot) RoamingMessages(timeStart, timeEnd, target int64) ([]any, error) {
	result, err := b.request("roamingMessages", "", &struct {
		TimeStart int64 `json:"timeStart"`
		TimeEnd   int64 `json:"timeEnd"`
		Target    int64 `json:"target"`
	}{timeStart, timeEnd, target})
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
			}
		}
		e := fmt.Sprint("decode message failed:", data.Raw)
		log.Errorln(e)
		return nil, errors.New(e)
	}
	return retArray, nil
}
