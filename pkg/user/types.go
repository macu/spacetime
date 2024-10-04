package user

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleModerator UserRole = "moderator"
	RoleUser      UserRole = "user"
	RoleInactive  UserRole = "inactive"
	RoleBanned    UserRole = "banned"
)

func CheckRoleActive(role string) bool {
	return role != string(RoleBanned) && role != string(RoleInactive)
}
