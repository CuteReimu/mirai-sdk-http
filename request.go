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

// MessageFromId 通过messageId获取消息
type MessageFromId struct {
	MessageId int64 `json:"messageId"` // 获取消息的messageId
	Target    int64 `json:"target"`    // 好友id或群id
}

func (m *MessageFromId) GetCommand() string {
	return "messageFromId"
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

// SendGroupMessage 获取群消息
type SendGroupMessage struct {
	Target       int64           `json:"target,omitempty"`
	Group        int64           `json:"group,omitempty"`
	Quote        int64           `json:"quote,omitempty"`
	MessageChain []SingleMessage `json:"messageChain"`
}

func (m *SendGroupMessage) GetCommand() string {
	return "sendGroupMessage"
}
