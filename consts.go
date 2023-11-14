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
