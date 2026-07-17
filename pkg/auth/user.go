package auth

import "database/sql"

func GetUserByHandle(db *sql.DB, handle string) (*User, error) {

	var user User

	err := db.QueryRow(`SELECT id, handle
		FROM user_account
		WHERE handle = $1
	`, handle).Scan(&user.ID, &user.Handle)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
