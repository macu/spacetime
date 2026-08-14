package user

import (
	"spacetime/pkg/utils/db"
)

func CheckAdmin(db db.DBConn, userID uint) bool {
	var userRole string
	err := db.QueryRow(`SELECT role FROM user_account WHERE id = $1`, userID).Scan(&userRole)
	if err != nil {
		return false
	}
	return userRole == string(RoleAdmin)
}
