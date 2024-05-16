package miraihttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"log/slog"
	"strconv"
)

type MessageChain []SingleMessage

func (c *MessageChain) UnmarshalJSON(data []byte) error {
	if !gjson.ValidBytes(data) {
		return errors.New("invalid json data")
	}
	result := gjson.ParseBytes(data)
	if !result.IsArray() {
		return errors.New("result is not array")
	}
	*c = parseMessageChain(result.Array())
	return nil
}

type SingleMessage interface {
	FillMessageType()
}

// buildMessageChain 自动填上每个元素的Type字段
func buildMessageChain(messages MessageChain) MessageChain {
	for _, m := range messages {
		m.FillMessageType()
	}
	return messages
}

// Source 永远为chain的第一个元素
type Source struct {
	Type string `json:"type"`
	Id   int64  `json:"id"`   // 消息的识别号，用于引用回复
	Time int64  `json:"time"` // 时间戳
}

func (m *Source) FillMessageType() {
	m.Type = "Source"
}

func (m *Source) String() string {
	return ""
}

// Quote 引用回复
type Quote struct {
	Type     string `json:"type"`
	Id       int64  `json:"id"`       // 被引用回复的原消息的messageId
	GroupId  int64  `json:"groupId"`  // 被引用回复的原消息所接收的群号，当为好友消息时为0
	SenderId int64  `json:"senderId"` // 被引用回复的原消息的发送者的QQ号
	TargetId int64  `json:"targetId"` // 被引用回复的原消息的接收者的QQ号（或群号）
	Origin   []any  `json:"origin"`   // 被引用回复的原消息的消息链对象
}

func (m *Quote) FillMessageType() {
	m.Type = "Quote"
}

func (m *Quote) String() string {
	return "[引用]"
}

// At @消息
type At struct {
	Type    string `json:"type"`
	Target  int64  `json:"target"`            // 群员QQ号
	Display string `json:"display,omitempty"` // At时显示的文字，发送消息时无效，自动使用群名片
}

func (m *At) FillMessageType() {
	m.Type = "At"
}

func (m *At) String() string {
	return "@" + strconv.FormatInt(m.Target, 10)
}

// AtAll @全体消息
type AtAll struct {
	Type string `json:"type"`
}

func (m *AtAll) FillMessageType() {
	m.Type = "AtAll"
}

func (m *AtAll) String() string {
	return "@全体成员"
}

// Face QQ表情
type Face struct {
	Type   string `json:"type"`
	FaceId int32  `json:"faceId,omitempty"` // QQ表情编号，可选，优先高于name
	Name   string `json:"name,omitempty"`   // QQ表情拼音，可选
}

func (m *Face) FillMessageType() {
	m.Type = "Face"
}

func (m *Face) String() string {
	switch {
	case m.Name != "":
		return "[" + m.Name + "]"
	case m.FaceId != 0:
		return fmt.Sprintf("[表情:%d]", m.FaceId)
	default:
		return "[表情]"
	}
}

// Plain 文字消息
type Plain struct {
	Type string `json:"type"`
	Text string `json:"text"` // 文字消息
}

func (m *Plain) FillMessageType() {
	m.Type = "Plain"
}

func (m *Plain) String() string {
	return m.Text
}

// Image 图片（参数优先级imageId > url > path > base64）
type Image struct {
	Type    string `json:"type"`
	ImageId string `json:"imageId,omitempty"` // 图片的imageId，群图片与好友图片格式不同。不为空时将忽略url属性
	Url     string `json:"url,omitempty"`     // 图片的URL，发送时可作网络图片的链接；接收时为腾讯图片服务器的链接，可用于图片下载
	Path    string `json:"path,omitempty"`    // 图片的路径，发送本地图片，路径相对于 JVM 工作路径（默认是当前路径，可通过 -Duser.dir=...指定），也可传入绝对路径。
	Base64  string `json:"base64,omitempty"`  // 图片的 Base64 编码
}

func (m *Image) FillMessageType() {
	m.Type = "Image"
}

func (m *Image) String() string {
	return "[图片]"
}

// FlashImage 闪照，参数同Image
type FlashImage struct {
	Type    string `json:"type"`
	ImageId string `json:"imageId,omitempty"`
	Url     string `json:"url,omitempty"`
	Path    string `json:"path,omitempty"`
	Base64  string `json:"base64,omitempty"`
}

func (m *FlashImage) FillMessageType() {
	m.Type = "FlashImage"
}

func (m *FlashImage) String() string {
	return "[闪照]"
}

// Voice 语音（参数优先级imageId > url > path > base64）
type Voice struct {
	Type    string `json:"type"`
	VoiceId string `json:"voiceId,omitempty"` // 语音的voiceId，不为空时将忽略url属性
	Url     string `json:"url,omitempty"`     // 语音的URL，发送时可作网络语音的链接；接收时为腾讯语音服务器的链接，可用于语音下载
	Path    string `json:"path,omitempty"`    // 语音的路径，发送本地语音，路径相对于 JVM 工作路径（默认是当前路径，可通过 -Duser.dir=...指定），也可传入绝对路径。
	Base64  string `json:"base64,omitempty"`  // 语音的 Base64 编码
	Length  string `json:"length,omitempty"`  // 返回的语音长度, 发送消息时可以不传
}

func (m *Voice) FillMessageType() {
	m.Type = "Voice"
}

func (m *Voice) String() string {
	return "[语音消息]"
}

type Xml struct {
	Type string `json:"type"`
	Xml  string `json:"xml"` // XML文本
}

func (m *Xml) FillMessageType() {
	m.Type = "Xml"
}

func (m *Xml) String() string {
	return m.Xml
}

type Json struct {
	Type string `json:"type"`
	Json string `json:"json"` // Json文本
}

func (m *Json) FillMessageType() {
	m.Type = "Json"
}

func (m *Json) String() string {
	return m.Json
}

type App struct {
	Type    string `json:"type"`
	Content string `json:"content"` // 内容
}

func (m *App) FillMessageType() {
	m.Type = "App"
}

func (m *App) String() string {
	return m.Content
}

// PokeName 戳一戳的类型
type PokeName string

const (
	PokeNamePoke        PokeName = "Poke"        // 戳一戳
	PokeNameShowLove    PokeName = "ShowLove"    // 比心
	PokeNameLike        PokeName = "Like"        // 点赞
	PokeNameHeartbroken PokeName = "Heartbroken" // 心碎
	PokeNameSixSixSix   PokeName = "SixSixSix"   // 666
	PokeNameFangDaZhao  PokeName = "FangDaZhao"  // 放大招
)

// Poke 戳一戳
type Poke struct {
	Type string   `json:"type"`
	Name PokeName `json:"name"` // 戳一戳的类型
}

func (m *Poke) FillMessageType() {
	m.Type = "Poke"
}

func (m *Poke) String() string {
	return "[戳一戳]"
}

// Dice 骰子
type Dice struct {
	Type  string `json:"type"`
	Value int32  `json:"value"` // 点数
}

func (m *Dice) FillMessageType() {
	m.Type = "Dice"
}

func (m *Dice) String() string {
	return fmt.Sprintf("[骰子:%d]", m.Value)
}

// MarketFace 商城表情（目前商城表情仅支持接收和转发，不支持构造发送）
type MarketFace struct {
	Type string `json:"type"`
	Id   int32  `json:"id"`   // 商城表情唯一标识
	Name string `json:"name"` // 表情显示名称
}

func (m *MarketFace) FillMessageType() {
	m.Type = "MarketFace"
}

func (m *MarketFace) String() string {
	switch {
	case m.Name != "":
		return "[" + m.Name + "]"
	case m.Id != 0:
		return fmt.Sprintf("[商城表情:%d]", m.Id)
	default:
		return "[商城表情]"
	}
}

// MusicShare 音乐分享
type MusicShare struct {
	Type       string `json:"type"`
	Kind       string `json:"kind"`       // 类型
	Title      string `json:"title"`      // 标题
	Summary    string `json:"summary"`    // 概括
	JumpUrl    string `json:"jumpUrl"`    // 跳转路径
	PictureUrl string `json:"pictureUrl"` // 封面路径
	MusicUrl   string `json:"musicUrl"`   // 音源路径
	Brief      string `json:"brief"`      // 简介
}

func (m *MusicShare) FillMessageType() {
	m.Type = "MusicShare"
}

func (m *MusicShare) String() string {
	return "[分享]" + m.Title
}

type ForwardMessageNode struct {
	SenderId     int64  `json:"senderId,omitempty"`     // 消息节点
	Time         int64  `json:"time,omitempty"`         // 发送时间
	SenderName   string `json:"senderName,omitempty"`   // 显示名称
	MessageChain []any  `json:"messageChain,omitempty"` // 消息数组

	MessageId int64 `json:"messageId,omitempty"` // 可以只使用消息messageId，从当前对话上下文缓存中读取一条消息作为节点

	// MessageRef 引用缓存中其他对话上下文的消息作为节点
	//
	// 参考 https://docs.mirai.mamoe.net/mirai-api-http/api/MessageType.html#forwardmessage
	MessageRef map[string]any `json:"messageRef,omitempty"`

	// (senderId, time, senderName, messageChain), messageId, messageRef 是三种不同构造引用节点的方式，选其中一个/组传参即可
}

// ForwardMessage 转发消息
type ForwardMessage struct {
	Type string `json:"type"`

	// Display 转发消息的卡片显示文本，值为表示使用客户端默认值。发送时可以直接填nil，表示全用默认值。
	//
	// 参考 https://docs.mirai.mamoe.net/mirai-api-http/api/MessageType.html#forwardmessage
	Display map[string]any `json:"display,omitempty"`

	NodeList []*ForwardMessageNode `json:"nodeList"` // 消息节点
}

func (m *ForwardMessage) FillMessageType() {
	m.Type = "Forward"
}

func (m *ForwardMessage) String() string {
	return "[转发消息]"
}

// File 文件
type File struct {
	Type string `json:"type"`
	Id   string `json:"id"`   // 文件识别id
	Name string `json:"name"` // 文件名
	Size int64  `json:"size"` // 文件大小
}

func (m *File) FillMessageType() {
	m.Type = "File"
}

func (m *File) String() string {
	return "[文件]" + m.Name
}

type MiraiCode struct {
	Type string `json:"type"`
	Code string `json:"code"`
}

func (m *MiraiCode) FillMessageType() {
	m.Type = "MiraiCode"
}

func (m *MiraiCode) String() string {
	return m.Code
}

var singleMessageBuilder = map[string]func() SingleMessage{
	"Source":     func() SingleMessage { return &Source{} },
	"Quote":      func() SingleMessage { return &Quote{} },
	"At":         func() SingleMessage { return &At{} },
	"AtAll":      func() SingleMessage { return &AtAll{} },
	"Face":       func() SingleMessage { return &Face{} },
	"Plain":      func() SingleMessage { return &Plain{} },
	"Image":      func() SingleMessage { return &Image{} },
	"FlashImage": func() SingleMessage { return &FlashImage{} },
	"Voice":      func() SingleMessage { return &Voice{} },
	"Xml":        func() SingleMessage { return &Xml{} },
	"Json":       func() SingleMessage { return &Json{} },
	"App":        func() SingleMessage { return &App{} },
	"Poke":       func() SingleMessage { return &Poke{} },
	"Dice":       func() SingleMessage { return &Dice{} },
	"MarketFace": func() SingleMessage { return &MarketFace{} },
	"MusicShare": func() SingleMessage { return &MusicShare{} },
	"Forward":    func() SingleMessage { return &ForwardMessage{} },
	"File":       func() SingleMessage { return &File{} },
	"MiraiCode":  func() SingleMessage { return &MiraiCode{} },
}

func parseMessageChain(results []gjson.Result) MessageChain {
	if len(results) == 0 {
		return nil
	}
	ret := make(MessageChain, 0, len(results))
	for i := range results {
		if results[i].Type != gjson.JSON {
			slog.Error("single message is not json: " + results[i].Type.String())
			continue
		}
		singleMessageType := results[i].Get("type").String()
		if builder, ok := singleMessageBuilder[singleMessageType]; ok {
			m := builder()
			if err := json.Unmarshal([]byte(results[i].Raw), m); err == nil {
				ret = append(ret, m)
			} else {
				slog.Error("json unmarshal failed", "buf", results[i].Raw, "error", err)
			}
		} else {
			slog.Error("unknown single message type: " + results[i].String())
		}
	}
	return ret
}
