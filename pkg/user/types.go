package user

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleModerator UserRole = "moderator"
	RoleUser      UserRole = "user"
	RoleInactive  UserRole = "inactive"
	RoleBanned    UserRole = "banned"
)

type Auth struct {
	UserID uint
	Role   UserRole
}

func CheckRoleValid(role string) bool {
	switch role {
	case string(RoleAdmin), string(RoleModerator), string(RoleUser),
		string(RoleInactive), string(RoleBanned):
		return true
	}
	return false
}

func CheckRoleActive(role string) bool {
	return !(role == string(RoleBanned) || role == string(RoleInactive))
}

func CheckRoleAdmin(role string) bool {
	return role == string(RoleAdmin)
}
