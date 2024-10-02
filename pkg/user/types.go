package user

type UserRole string

const (
	RoleAdmin       UserRole = "admin"
	RoleModerator   UserRole = "moderator"
	RoleContributor UserRole = "contributor"
	RoleBanned      UserRole = "banned"
)

func CheckBanned(role string) bool {
	return role == string(RoleBanned)
}
