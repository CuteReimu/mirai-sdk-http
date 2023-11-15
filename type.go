package miraihttp

type Perm string

const (
	PermOwner         Perm = "OWNER"         // 群主
	PermAdministrator Perm = "ADMINISTRATOR" // 管理员
	PermMember        Perm = "MEMBER"        // 群成员
)

type Kind string

const (
	KindFriend   Kind = "Friend"
	KindGroup    Kind = "Group"
	KindStranger Kind = "Stranger"
)

// Friend 好友
type Friend struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
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
