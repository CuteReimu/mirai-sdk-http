package miraihttp

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

func init() {
	decoder["NewFriendRequestEvent"] = parseEvent[NewFriendRequestEvent]
	decoder["MemberJoinRequestEvent"] = parseEvent[MemberJoinRequestEvent]
	decoder["BotInvitedJoinGroupRequestEvent"] = parseEvent[BotInvitedJoinGroupRequestEvent]
	decoder["BotGroupPermissionChangeEvent"] = parseEvent[BotGroupPermissionChangeEvent]
	decoder["BotMuteEvent"] = parseEvent[BotMuteEvent]
	decoder["BotUnmuteEvent"] = parseEvent[BotUnmuteEvent]
	decoder["BotJoinGroupEvent"] = parseEvent[BotJoinGroupEvent]
	decoder["BotLeaveEventActive"] = parseEvent[BotLeaveEventActive]
	decoder["BotLeaveEventKick"] = parseEvent[BotLeaveEventKick]
	decoder["BotLeaveEventDisband"] = parseEvent[BotLeaveEventDisband]
	decoder["GroupRecallEvent"] = parseEvent[GroupRecallEvent]
	decoder["FriendRecallEvent"] = parseEvent[FriendRecallEvent]
	decoder["NudgeEvent"] = parseEvent[NudgeEvent]
	decoder["GroupNameChangeEvent"] = parseEvent[GroupNameChangeEvent]
	decoder["GroupEntranceAnnouncementChangeEvent"] = parseEvent[GroupEntranceAnnouncementChangeEvent]
	decoder["GroupMuteAllEvent"] = parseEvent[GroupMuteAllEvent]
	decoder["GroupAllowAnonymousChatEvent"] = parseEvent[GroupAllowAnonymousChatEvent]
	decoder["GroupAllowConfessTalkEvent"] = parseEvent[GroupAllowConfessTalkEvent]
	decoder["GroupAllowMemberInviteEvent"] = parseEvent[GroupAllowMemberInviteEvent]
	decoder["MemberJoinEvent"] = parseEvent[MemberJoinEvent]
	decoder["MemberLeaveEventKick"] = parseEvent[MemberLeaveEventKick]
	decoder["MemberLeaveEventQuit"] = parseEvent[MemberLeaveEventQuit]
	decoder["MemberCardChangeEvent"] = parseEvent[MemberCardChangeEvent]
	decoder["MemberSpecialTitleChangeEvent"] = parseEvent[MemberSpecialTitleChangeEvent]
	decoder["MemberPermissionChangeEvent"] = parseEvent[MemberPermissionChangeEvent]
	decoder["MemberMuteEvent"] = parseEvent[MemberMuteEvent]
	decoder["MemberUnmuteEvent"] = parseEvent[MemberUnmuteEvent]
	decoder["MemberHonorChangeEvent"] = parseEvent[MemberHonorChangeEvent]
}

func parseEvent[T any](data gjson.Result) any {
	var m T
	if err := json.Unmarshal([]byte(data.Raw), &m); err != nil {
		log.Errorln("json unmarshal failed: ", err)
		return nil
	}
	return &m
}

// BotGroupPermissionChangeEvent Bot在群里的权限被改变. 操作人一定是群主
type BotGroupPermissionChangeEvent struct {
	Origin  Perm  `json:"origin"`  // Bot的原权限
	Current Perm  `json:"current"` // Bot的新权限
	Group   Group `json:"group"`
}

// ListenBotGroupPermissionChangeEvent 监听Bot在群里的权限被改变
func (b *Bot) ListenBotGroupPermissionChangeEvent(l func(message *BotGroupPermissionChangeEvent) bool) {
	listen(b, "BotGroupPermissionChangeEvent", l)
}

// BotMuteEvent Bot被禁言
type BotMuteEvent struct {
	DurationSeconds int64  `json:"durationSeconds"` // 禁言时长，单位为秒
	Operator        Member `json:"operator"`
}

// ListenBotMuteEvent 监听Bot被禁言
func (b *Bot) ListenBotMuteEvent(l func(message *BotMuteEvent) bool) {
	listen(b, "BotMuteEvent", l)
}

// BotUnmuteEvent Bot被取消禁言
type BotUnmuteEvent struct {
	Operator Member `json:"operator"`
}

// ListenBotUnmuteEvent 监听Bot被取消禁言
func (b *Bot) ListenBotUnmuteEvent(l func(message *BotUnmuteEvent) bool) {
	listen(b, "BotUnmuteEvent", l)
}

// BotJoinGroupEvent Bot加入了一个新群
type BotJoinGroupEvent struct {
	Group   Group   `json:"group"`
	Invitor *Member `json:"invitor,omitempty"` // 邀请者，可能为空
}

// ListenBotJoinGroupEvent 监听Bot加入了一个新群
func (b *Bot) ListenBotJoinGroupEvent(l func(message *BotJoinGroupEvent) bool) {
	listen(b, "BotJoinGroupEvent", l)
}

// BotLeaveEventActive Bot主动退出一个群
type BotLeaveEventActive struct {
	Group Group `json:"group"`
}

// ListenBotLeaveEventActive 监听Bot主动退出一个群
func (b *Bot) ListenBotLeaveEventActive(l func(message *BotLeaveEventActive) bool) {
	listen(b, "BotLeaveEventActive", l)
}

// BotLeaveEventKick Bot被踢出一个群
type BotLeaveEventKick struct {
	Group Group `json:"group"`
}

// ListenBotLeaveEventKick 监听Bot被踢出一个群
func (b *Bot) ListenBotLeaveEventKick(l func(message *BotLeaveEventKick) bool) {
	listen(b, "BotLeaveEventKick", l)
}

// BotLeaveEventDisband Bot因群主解散群而退出群, 操作人一定是群主
type BotLeaveEventDisband struct {
	Group    Group   `json:"group"`
	Operator *Member `json:"operator"`
}

// ListenBotLeaveEventDisband 监听Bot因群主解散群而退出群
func (b *Bot) ListenBotLeaveEventDisband(l func(message *BotLeaveEventDisband) bool) {
	listen(b, "BotLeaveEventDisband", l)
}

// GroupRecallEvent 群消息撤回
type GroupRecallEvent struct {
	AuthorId  int64  `json:"authorId"`  // 原消息发送者的QQ号
	MessageId int64  `json:"messageId"` // 原消息messageId
	Time      int64  `json:"time"`      // 原消息发送时间
	Group     Group  `json:"group"`     // 消息撤回所在的群
	Operator  Member `json:"operator"`  // 撤回消息的操作人，当null时为bot操作
}

// ListenGroupRecallEvent 监听群消息撤回
func (b *Bot) ListenGroupRecallEvent(l func(message *GroupRecallEvent) bool) {
	listen(b, "GroupRecallEvent", l)
}

// FriendRecallEvent 好友消息撤回
type FriendRecallEvent struct {
	AuthorId  int64 `json:"authorId"`  // 原消息发送者的QQ号
	MessageId int64 `json:"messageId"` // 原消息messageId
	Time      int64 `json:"time"`      // 原消息发送时间
	Operator  int64 `json:"operator"`  // 好友QQ号或BotQQ号
}

// ListenFriendRecallEvent 监听好友消息撤回
func (b *Bot) ListenFriendRecallEvent(l func(message *FriendRecallEvent) bool) {
	listen(b, "FriendRecallEvent", l)
}

// NudgeEvent 戳一戳事件
type NudgeEvent struct {
	FromId  int64 `json:"fromId"` // 动作发出者的QQ号
	Subject struct {
		Id   int64 `json:"id"`   // 来源的QQ号（好友）或群号
		Kind Kind  `json:"kind"` // 来源的类型
	} `json:"subject"`
	Action string `json:"action"` // 动作类型
	Suffix string `json:"suffix"` // 自定义动作内容
	Target int64  `json:"target"` // 动作目标的QQ号
}

// ListenNudgeEvent 监听戳一戳事件
func (b *Bot) ListenNudgeEvent(l func(message *NudgeEvent) bool) {
	listen(b, "NudgeEvent", l)
}

// GroupNameChangeEvent 某个群名改变
type GroupNameChangeEvent struct {
	Origin   string `json:"origin"`  // 原群名
	Current  string `json:"current"` // 新群名
	Group    Group  `json:"group"`
	Operator Member `json:"operator"`
}

// ListenGroupNameChangeEvent 监听某个群名改变
func (b *Bot) ListenGroupNameChangeEvent(l func(message *GroupNameChangeEvent) bool) {
	listen(b, "GroupNameChangeEvent", l)
}

// GroupEntranceAnnouncementChangeEvent 某群入群公告改变
type GroupEntranceAnnouncementChangeEvent struct {
	Origin   string `json:"origin"`  // 原公告
	Current  string `json:"current"` // 新公告
	Group    Group  `json:"group"`
	Operator Member `json:"operator"`
}

// ListenGroupEntranceAnnouncementChangeEvent 监听某群入群公告改变
func (b *Bot) ListenGroupEntranceAnnouncementChangeEvent(l func(message *GroupEntranceAnnouncementChangeEvent) bool) {
	listen(b, "GroupEntranceAnnouncementChangeEvent", l)
}

// GroupMuteAllEvent 全员禁言
type GroupMuteAllEvent struct {
	Origin   bool   `json:"origin"`
	Current  bool   `json:"current"`
	Group    Group  `json:"group"`
	Operator Member `json:"operator"`
}

// ListenGroupMuteAllEvent 监听全员禁言
func (b *Bot) ListenGroupMuteAllEvent(l func(message *GroupMuteAllEvent) bool) {
	listen(b, "GroupMuteAllEvent", l)
}

// GroupAllowAnonymousChatEvent 匿名聊天
type GroupAllowAnonymousChatEvent struct {
	Origin   bool   `json:"origin"`
	Current  bool   `json:"current"`
	Group    Group  `json:"group"`
	Operator Member `json:"operator"`
}

// ListenGroupAllowAnonymousChatEvent 监听匿名聊天
func (b *Bot) ListenGroupAllowAnonymousChatEvent(l func(message *GroupAllowAnonymousChatEvent) bool) {
	listen(b, "GroupAllowAnonymousChatEvent", l)
}

// GroupAllowConfessTalkEvent 坦白说
type GroupAllowConfessTalkEvent struct {
	Origin  bool  `json:"origin"`
	Current bool  `json:"current"`
	Group   Group `json:"group"`
	IsByBot bool  `json:"isByBot"` // 是否Bot进行该操作
}

// ListenGroupAllowConfessTalkEvent 监听坦白说
func (b *Bot) ListenGroupAllowConfessTalkEvent(l func(message *GroupAllowConfessTalkEvent) bool) {
	listen(b, "GroupAllowConfessTalkEvent", l)
}

// GroupAllowMemberInviteEvent 允许群员邀请好友加群
type GroupAllowMemberInviteEvent struct {
	Origin   bool   `json:"origin"`
	Current  bool   `json:"current"`
	Group    Group  `json:"group"`
	Operator Member `json:"operator"`
}

// ListenGroupAllowMemberInviteEvent 监听允许群员邀请好友加群
func (b *Bot) ListenGroupAllowMemberInviteEvent(l func(message *GroupAllowMemberInviteEvent) bool) {
	listen(b, "GroupAllowMemberInviteEvent", l)
}

// MemberJoinEvent 新人入群的事件
type MemberJoinEvent struct {
	Member  Member  `json:"member"`
	Invitor *Member `json:"invitor"` // 邀请人，可能为空
}

// ListenMemberJoinEvent 监听新人入群的事件
func (b *Bot) ListenMemberJoinEvent(l func(message *MemberJoinEvent) bool) {
	listen(b, "MemberJoinEvent", l)
}

// MemberLeaveEventKick 成员被踢出群（该成员不是Bot）
type MemberLeaveEventKick struct {
	Member   Member `json:"member"`
	Operator Member `json:"operator"`
}

// ListenMemberLeaveEventKick 监听成员被踢出群
func (b *Bot) ListenMemberLeaveEventKick(l func(message *MemberLeaveEventKick) bool) {
	listen(b, "MemberLeaveEventKick", l)
}

// MemberLeaveEventQuit 成员主动离群（该成员不是Bot）
type MemberLeaveEventQuit struct {
	Member Member `json:"member"`
}

// ListenMemberLeaveEventQuit 监听成员主动离群
func (b *Bot) ListenMemberLeaveEventQuit(l func(message *MemberLeaveEventQuit) bool) {
	listen(b, "MemberLeaveEventQuit", l)
}

// MemberCardChangeEvent 群名片改动
type MemberCardChangeEvent struct {
	Origin  string `json:"origin"`
	Current string `json:"current"`
	Member  Member `json:"member"` // 名片改动的群员的信息
}

// ListenMemberCardChangeEvent 监听群名片改动
func (b *Bot) ListenMemberCardChangeEvent(l func(message *MemberCardChangeEvent) bool) {
	listen(b, "MemberCardChangeEvent", l)
}

// MemberSpecialTitleChangeEvent 群头衔改动（只有群主有操作限权）
type MemberSpecialTitleChangeEvent struct {
	Origin  string `json:"origin"`
	Current string `json:"current"`
	Member  Member `json:"member"`
}

// ListenMemberSpecialTitleChangeEvent 监听群头衔改动
func (b *Bot) ListenMemberSpecialTitleChangeEvent(l func(message *MemberSpecialTitleChangeEvent) bool) {
	listen(b, "MemberSpecialTitleChangeEvent", l)
}

// MemberPermissionChangeEvent 成员权限改变的事件（该成员不是Bot）
type MemberPermissionChangeEvent struct {
	Origin  Perm   `json:"origin"`
	Current Perm   `json:"current"`
	Member  Member `json:"member"`
}

// ListenMemberPermissionChangeEvent 监听成员权限改变的事件
func (b *Bot) ListenMemberPermissionChangeEvent(l func(message *MemberPermissionChangeEvent) bool) {
	listen(b, "MemberPermissionChangeEvent", l)
}

// MemberMuteEvent 群成员被禁言事件（该成员不是Bot）
type MemberMuteEvent struct {
	DurationSeconds int64  `json:"durationSeconds"`
	Member          Member `json:"member"`
	Operator        Member `json:"operator"`
}

// ListenMemberMuteEvent 监听群成员被禁言事件
func (b *Bot) ListenMemberMuteEvent(l func(message *MemberMuteEvent) bool) {
	listen(b, "MemberMuteEvent", l)
}

// MemberUnmuteEvent 群成员被取消禁言事件（该成员不是Bot）
type MemberUnmuteEvent struct {
	Member   Member `json:"member"`
	Operator Member `json:"operator"`
}

// ListenMemberUnmuteEvent 监听群成员被取消禁言事件
func (b *Bot) ListenMemberUnmuteEvent(l func(message *MemberUnmuteEvent) bool) {
	listen(b, "MemberUnmuteEvent", l)
}

// MemberHonorChangeEvent 群员称号改变
type MemberHonorChangeEvent struct {
	Member Member      `json:"member"`
	Action HonorAction `json:"action"` // 称号变化行为
	Honor  string      `json:"honor"`  // 称号名称
}

// ListenMemberHonorChangeEvent 监听群员称号改变
func (b *Bot) ListenMemberHonorChangeEvent(l func(message *MemberHonorChangeEvent) bool) {
	listen(b, "MemberHonorChangeEvent", l)
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
