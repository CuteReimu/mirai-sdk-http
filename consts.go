package miraihttp

type Perm string

const (
	PermOwner         Perm = "OWNER"         // 群主
	PermAdministrator Perm = "ADMINISTRATOR" // 管理员
	PermMember        Perm = "MEMBER"        // 群成员
)
