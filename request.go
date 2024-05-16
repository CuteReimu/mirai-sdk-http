package miraihttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"log/slog"
)

// About 获取插件版本号
func (b *Bot) About() (string, error) {
	result, err := b.request2("about", "", nil)
	if err != nil {
		return "", err
	}
	return result.Get("data.version").String(), nil
}

// BotList 获取登录账号
func (b *Bot) BotList() ([]int64, error) {
	result, err := b.request2("botList", "", nil)
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
	result, err := b.request2("messageFromId", "", &struct {
		MessageId int64 `json:"messageId"`
		Target    int64 `json:"target"`
	}{messageId, target})
	if err != nil {
		return nil, err
	}
	data := result.Get("data")
	if data.Type != gjson.JSON {
		e := fmt.Sprint("invalid json message: ", result)
		slog.Error(e)
		return nil, errors.New(e)
	}
	messageType := data.Get("type").String()
	if p := decoder[messageType]; p != nil {
		if m := p(data); m != nil {
			return m, nil
		}
	}
	e := fmt.Sprint("decode message failed:", data.Raw)
	slog.Error(e)
	return nil, errors.New(e)
}

// SendFriendMessage 发送好友消息，qq-目标好友的QQ号，quote-引用回复的消息，messageChain-发送的内容，返回消息id
func (b *Bot) SendFriendMessage(qq, quote int64, messageChain MessageChain) (int64, error) {
	result, err := b.request2("sendFriendMessage", "", &struct {
		Target       int64        `json:"target"`
		Quote        int64        `json:"quote,omitempty"`
		MessageChain MessageChain `json:"messageChain"`
	}{qq, quote, buildMessageChain(messageChain)}, "messageId")
	if err != nil {
		return 0, err
	}
	return result.Int(), nil
}

// SendGroupMessage 发送群消息，group-群号，quote-引用回复的消息，messageChain-发送的内容，返回消息id
func (b *Bot) SendGroupMessage(group, quote int64, messageChain MessageChain) (int64, error) {
	result, err := b.request2("sendGroupMessage", "", &struct {
		Target       int64        `json:"target"`
		Quote        int64        `json:"quote,omitempty"`
		MessageChain MessageChain `json:"messageChain"`
	}{group, quote, buildMessageChain(messageChain)}, "messageId")
	if err != nil {
		return 0, err
	}
	return result.Int(), nil
}

// SendTempMessage 发送临时会话消息，qq-临时会话对象QQ号，group-临时会话群号，quote-引用回复的消息，messageChain-发送的内容，返回消息id
func (b *Bot) SendTempMessage(qq, group, quote int64, messageChain MessageChain) (int64, error) {
	result, err := b.request2("sendTempMessage", "", &struct {
		QQ           int64        `json:"qq"`
		Group        int64        `json:"group"`
		Quote        int64        `json:"quote,omitempty"`
		MessageChain MessageChain `json:"messageChain"`
	}{qq, group, quote, buildMessageChain(messageChain)}, "messageId")
	if err != nil {
		return 0, err
	}
	return result.Int(), nil
}

// SendNudge 发送头像戳一戳消息，qq-戳谁，subject-这条消息发到哪（好友/群），kind-上下文类型
func (b *Bot) SendNudge(qq, subject int64, kind Kind) error {
	_, err := b.request2("sendNudge", "", &struct {
		Target  int64 `json:"target"`
		Subject int64 `json:"subject"`
		Kind    Kind  `json:"kind"`
	}{qq, subject, kind})
	return err
}

// Recall 撤回消息，target-撤回哪的消息（好友/群），messageId-需要撤回的消息的messageId
func (b *Bot) Recall(target, messageId int64) error {
	_, err := b.request2("recall", "", &struct {
		Target    int64 `json:"target"`
		MessageId int64 `json:"messageId"`
	}{target, messageId})
	return err
}

// RoamingMessages 获取漫游消息，timeStart和timeEnd为开始和结束的时间戳，单位为秒。qq为查询的对象QQ，目前仅支持好友漫游消息。
//
// 返回数组的元素为 FriendMessage, GroupMessage, TempMessage, StrangerMessage
func (b *Bot) RoamingMessages(timeStart, timeEnd, qq int64) ([]any, error) {
	result, err := b.request2("roamingMessages", "", &struct {
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
			slog.Error(e)
			return nil, errors.New(e)
		}
		messageType := data.Get("type").String()
		if p := decoder[messageType]; p != nil {
			if m := p(data); m != nil {
				retArray = append(retArray, m)
				continue
			}
		}
		e := fmt.Sprint("decode message failed:", data.Raw)
		slog.Error(e)
		return nil, errors.New(e)
	}
	return retArray, nil
}

// Mute 禁言群成员（需要有相关限权），group-群，qq-被禁言的人，time-时间，单位秒，最多30天
func (b *Bot) Mute(group, qq, time int64) error {
	_, err := b.request2("mute", "", &struct {
		Target   int64 `json:"target"`
		MemberId int64 `json:"memberId"`
		Time     int64 `json:"time"`
	}{group, qq, time})
	return err
}

// Unmute 解除禁言群成员（需要有相关限权），group-群，qq-解除禁言的人
func (b *Bot) Unmute(group, qq int64) error {
	_, err := b.request2("unmute", "", &struct {
		Target   int64 `json:"target"`
		MemberId int64 `json:"memberId"`
	}{group, qq})
	return err
}

// Kick 移除群成员（需要有相关限权），group-群，qq-移除的人，block-移除后是否拉黑，msg-信息
func (b *Bot) Kick(group, qq int64, block bool, msg string) error {
	_, err := b.request2("kick", "", &struct {
		Target   int64  `json:"target"`
		MemberId int64  `json:"memberId"`
		Block    bool   `json:"block"`
		Msg      string `json:"msg"`
	}{group, qq, block, msg})
	return err
}

// Quit 退出群聊（自己不能是群主）
func (b *Bot) Quit(group int64) error {
	_, err := b.request2("quit", "", &struct {
		Target int64 `json:"target"`
	}{group})
	return err
}

// MuteAll 全体禁言（需要有相关限权）
func (b *Bot) MuteAll(group int64) error {
	_, err := b.request2("muteAll", "", &struct {
		Target int64 `json:"target"`
	}{group})
	return err
}

// UnmuteAll 解除全体禁言（需要有相关限权）
func (b *Bot) UnmuteAll(group int64) error {
	_, err := b.request2("unmuteAll", "", &struct {
		Target int64 `json:"target"`
	}{group})
	return err
}

// SetEssence 设置群精华消息（需要有相关限权）
func (b *Bot) SetEssence(group, messageId int64) error {
	_, err := b.request2("setEssence", "", &struct {
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
		slog.Error(e)
		return nil, err
	}
	return groupConfig, nil
}

// UpdateGroupConfig 修改群设置（需要有相关限权）
func (b *Bot) UpdateGroupConfig(group int64, groupConfig *GroupConfig) error {
	_, err := b.request2("groupConfig", "update", &struct {
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
		slog.Error(e)
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
	_, err := b.request2("memberInfo", "update", &struct {
		Target   int64 `json:"target"`
		MemberId int64 `json:"memberId"`
		Info     Info  `json:"info"`
	}{group, qq, Info{name, specialTitle}})
	return err
}

// MemberAdmin 修改群员管理员（需要有群主限权），assign-是否设置为管理员
func (b *Bot) MemberAdmin(group, qq int64, assign bool) error {
	_, err := b.request2("memberAdmin", "", &struct {
		Target   int64 `json:"target"`
		MemberId int64 `json:"memberId"`
		Assign   bool  `json:"assign"`
	}{group, qq, assign})
	return err
}

// FriendList 获取好友列表
func (b *Bot) FriendList() ([]*Friend, error) {
	result, err := b.request2("friendList", "", nil)
	if err != nil {
		return nil, err
	}
	var friends []*Friend
	if err = json.Unmarshal([]byte(result.Raw), &friends); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return friends, nil
}

// GroupList 获取群列表
func (b *Bot) GroupList() ([]*Group, error) {
	result, err := b.request2("groupList", "", nil)
	if err != nil {
		return nil, err
	}
	var groups []*Group
	if err = json.Unmarshal([]byte(result.Raw), &groups); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return groups, nil
}

// MemberList 获取群成员列表
func (b *Bot) MemberList(group int64) ([]*Member, error) {
	result, err := b.request2("memberList", "", &struct {
		Target int64 `json:"target"`
	}{group})
	if err != nil {
		return nil, err
	}
	var members []*Member
	if err = json.Unmarshal([]byte(result.Raw), &members); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return members, nil
}

// LatestMemberList 获取最新群成员列表，qqs为空表示获取所有
func (b *Bot) LatestMemberList(group int64, qqs []int64) ([]*Member, error) {
	result, err := b.request2("latestMemberList", "", &struct {
		Target    int64   `json:"target"`
		MemberIds []int64 `json:"memberIds"`
	}{group, qqs})
	if err != nil {
		return nil, err
	}
	var members []*Member
	if err = json.Unmarshal([]byte(result.Raw), &members); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return members, nil
}

// BotProfile 获取Bot资料
func (b *Bot) BotProfile() (*Profile, error) {
	result, err := b.request("botProfile", "", nil)
	if err != nil {
		return nil, err
	}
	profile := &Profile{}
	if err = json.Unmarshal([]byte(result.Raw), profile); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return profile, nil
}

// FriendProfile 获取好友资料
func (b *Bot) FriendProfile(qq int64) (*Profile, error) {
	result, err := b.request("friendProfile", "", &struct {
		Target int64 `json:"target"`
	}{qq})
	if err != nil {
		return nil, err
	}
	profile := &Profile{}
	if err = json.Unmarshal([]byte(result.Raw), profile); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return profile, nil
}

// MemberProfile 获取群成员资料
func (b *Bot) MemberProfile(group, qq int64) (*Profile, error) {
	result, err := b.request("memberProfile", "", &struct {
		Target   int64 `json:"target"`
		MemberId int64 `json:"memberId"`
	}{group, qq})
	if err != nil {
		return nil, err
	}
	profile := &Profile{}
	if err = json.Unmarshal([]byte(result.Raw), profile); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return profile, nil
}

// UserProfile 获取QQ用户资料
func (b *Bot) UserProfile(qq int64) (*Profile, error) {
	result, err := b.request("userProfile", "", &struct {
		Target int64 `json:"target"`
	}{qq})
	if err != nil {
		return nil, err
	}
	profile := &Profile{}
	if err = json.Unmarshal([]byte(result.Raw), profile); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return profile, nil
}

type FileParam struct {
	Id     string `json:"id"`               // 文件夹id, 空串为根目录
	Path   string `json:"path,omitempty"`   // 文件夹路径, 文件夹允许重名, 不保证准确, 准确定位使用 id
	Target int64  `json:"target,omitempty"` // 群号或好友QQ号
	Group  int64  `json:"group,omitempty"`  // 群号
	QQ     int64  `json:"qq,omitempty"`     // 好友QQ号

	// 以下是查看文件列表 GetFileList 获取文件信息 GetFileInfo 时可选

	WithDownloadInfo bool `json:"withDownloadInfo,omitempty"` // 是否携带下载信息。额外请求，无必要不要携带

	// 以下是查看文件列表 GetFileList 时需要

	Offset int `json:"offset,omitempty"` // 分页偏移
	Size   int `json:"size,omitempty"`   // 分页大小

	// 以下是创建文件夹 FileMkdir 时需要

	DirectoryName string `json:"directoryName,omitempty"` // 新建文件夹名

	// 以下是移动文件 FileMove 时需要

	MoveTo     string `json:"moveTo,omitempty"`     // 移动目标文件夹id
	MoveToPath string `json:"moveToPath,omitempty"` // 移动目标文件路径, 文件夹允许重名, 不保证准确, 准确定位使用 MoveTo

	// 以下是重命名文件 FileRename 时需要

	RenameTo string `json:"renameTo,omitempty"` // 新文件名
}

type FileDownloadInfo struct {
	Sha1           string `json:"sha1"`
	Md5            string `json:"md5"`
	DownloadTimes  int    `json:"downloadTimes"`
	UploaderId     int    `json:"uploaderId"`
	UploadTime     int    `json:"uploadTime"`
	LastModifyTime int    `json:"lastModifyTime"`
	Url            string `json:"url"`
}

type FileInfo struct {
	Name         string `json:"name"`
	Id           string `json:"id"`
	Path         string `json:"path"`
	Parent       any    `json:"parent"`
	Contact      Group  `json:"contact"`
	IsFile       bool   `json:"isFile"`
	IsDictionary bool   `json:"isDictionary"`
	IsDirectory  bool   `json:"isDirectory"`

	// 以下字段只有查看文件列表 GetFileList 获取文件信息 GetFileInfo 时才会有

	Sha1           string            `json:"sha1"`
	Md5            string            `json:"md5"`
	DownloadTimes  int               `json:"downloadTimes"`
	UploaderId     int               `json:"uploaderId"`
	UploadTime     int               `json:"uploadTime"`
	LastModifyTime int               `json:"lastModifyTime"`
	DownloadInfo   *FileDownloadInfo `json:"downloadInfo"` // 只有 WithDownloadInfo 为 true 时才会有
}

// GetFileList 查看文件列表
func (b *Bot) GetFileList(param FileParam) ([]*FileInfo, error) {
	result, err := b.request2("file_list", "", param)
	if err != nil {
		return nil, err
	}
	var fileList []*FileInfo
	if err = json.Unmarshal([]byte(result.Raw), &fileList); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return fileList, nil
}

// GetFileInfo 获取文件信息
func (b *Bot) GetFileInfo(param FileParam) (*FileInfo, error) {
	result, err := b.request2("file_info", "", param)
	if err != nil {
		return nil, err
	}
	fileList := &FileInfo{}
	if err = json.Unmarshal([]byte(result.Raw), fileList); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return fileList, nil
}

// FileMkdir 创建文件夹
func (b *Bot) FileMkdir(param FileParam) (*FileInfo, error) {
	result, err := b.request2("file_mkdir", "", param)
	if err != nil {
		return nil, err
	}
	fileList := &FileInfo{}
	if err = json.Unmarshal([]byte(result.Raw), fileList); err != nil {
		e := fmt.Sprint("unmarshal json failed: ", err)
		slog.Error(e)
		return nil, err
	}
	return fileList, nil
}

// FileDelete 删除文件
func (b *Bot) FileDelete(param FileParam) error {
	result, err := b.request2("file_mkdir", "", param)
	if err != nil {
		return err
	}
	code := result.Get("code").Int()
	if code != 0 {
		e := fmt.Sprint("file_mkdir failed: ", result.Get("msg").String(), ", code:", code)
		slog.Error(e)
		return errors.New(e)
	}
	return nil
}

// FileMove 移动文件
func (b *Bot) FileMove(param FileParam) error {
	_, err := b.request2("file_move", "", param)
	return err
}

// FileRename 重命名文件
func (b *Bot) FileRename(param FileParam) error {
	_, err := b.request2("file_rename", "", param)
	return err
}

// ResponseNewFriend 处理添加好友申请。operate：0-同意，1-拒绝，2-拒绝并拉黑
func (b *Bot) ResponseNewFriend(request *NewFriendRequestEvent, operate int, message string) error {
	_, err := b.request("resp_newFriendRequestEvent", "", &struct {
		EventId int64  `json:"eventId"`
		FromId  int64  `json:"fromId"`
		GroupId int64  `json:"groupId"`
		Operate int    `json:"operate"`
		Message string `json:"message"`
	}{request.EventId, request.QQ, request.Group, operate, message})
	return err
}

// ResponseMemberJoin 处理用户入群申请，Bot需要有管理员权限。operate：0-同意，1-拒绝，2-忽略，3-拒绝并拉黑，4-忽略并拉黑
func (b *Bot) ResponseMemberJoin(request *MemberJoinRequestEvent, operate int, message string) error {
	_, err := b.request("resp_memberJoinRequestEvent", "", &struct {
		EventId int64  `json:"eventId"`
		FromId  int64  `json:"fromId"`
		GroupId int64  `json:"groupId"`
		Operate int    `json:"operate"`
		Message string `json:"message"`
	}{request.EventId, request.QQ, request.Group, operate, message})
	return err
}

// ResponseBotInvitedJoinGroup 处理Bot被邀请入群申请，operate：0-同意，1-拒绝
func (b *Bot) ResponseBotInvitedJoinGroup(request *BotInvitedJoinGroupRequestEvent, operate int, message string) error {
	_, err := b.request("resp_newFriendRequestEvent", "", &struct {
		EventId int64  `json:"eventId"`
		FromId  int64  `json:"fromId"`
		GroupId int64  `json:"groupId"`
		Operate int    `json:"operate"`
		Message string `json:"message"`
	}{request.EventId, request.QQ, request.Group, operate, message})
	return err
}
