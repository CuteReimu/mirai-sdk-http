package miraihttp

import "fmt"

type Perm string

const (
	PermOwner         Perm = "OWNER"         // 群主
	PermAdministrator Perm = "ADMINISTRATOR" // 管理员
	PermMember        Perm = "MEMBER"        // 群成员
)

type Kind string

const (
	KindFriend   Kind = "Friend"   // 好友
	KindGroup    Kind = "Group"    // 群
	KindStranger Kind = "Stranger" // 陌生人
)

type HonorAction string

const (
	HonorActionAchieve HonorAction = "achieve" // 获得称号
	HonorActionLose    HonorAction = "lose"    // 失去称号
)

// Friend 好友
type Friend struct {
	Id       int64  `json:"id"`       // QQ号
	Nickname string `json:"nickname"` // 昵称
	Remark   string `json:"remark"`   // 备注
}

func (f *Friend) String() string {
	return fmt.Sprintf("%s(%d)", f.Nickname, f.Id)
}

// Group 群
type Group struct {
	Id         int64  `json:"id"`         // 群号
	Name       string `json:"name"`       // 群名称
	Permission Perm   `json:"permission"` // Bot在群中的权限
}

func (g *Group) String() string {
	return fmt.Sprintf("%s(%d)", g.Name, g.Id)
}

// Member 群成员
type Member struct {
	Id                 int64  `json:"id"`                 // QQ号
	MemberName         string `json:"memberName"`         // 群名片
	SpecialTitle       string `json:"specialTitle"`       // 群头衔
	Permission         Perm   `json:"permission"`         // 在群中的权限
	JoinTimestamp      int64  `json:"joinTimestamp"`      // 入群时间时间戳
	LastSpeakTimestamp int64  `json:"lastSpeakTimestamp"` // 最近发言时间戳
	MuteTimeRemaining  int64  `json:"muteTimeRemaining"`  // 剩余禁言时长
	Group              Group  `json:"group"`              // 群信息
}

func (m *Member) String() string {
	return fmt.Sprintf("%s(%d)", m.MemberName, m.Id)
}

type Sex string

const (
	SexUnknown Sex = "UNKNOWN" // 未知
	SexMale    Sex = "MALE"    // 男
	SexFemale  Sex = "FEMALE"  // 女
)

// Profile 用户资料
type Profile struct {
	Nickname string `json:"nickname"` // 昵称
	Email    string `json:"email"`    // 邮箱
	Age      int    `json:"age"`      // 年龄
	Level    int    `json:"level"`    // 等级
	Sign     string `json:"sign"`     // 个性签名
	Sex      Sex    `json:"sex"`      // 性别
}
