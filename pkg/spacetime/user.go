package spacetime

import "database/sql"

func GetUserSpaceID(db *sql.DB, userID uint) (uint, error) {

	var spaceID uint

	err := db.QueryRow(`SELECT id
		FROM space
		INNER JOIN user_space ON space.id = user_space.space_id
		WHERE space.created_by = $1 AND user_space.user_id = $1
	`, userID).Scan(&spaceID)

	if err != nil {
		return 0, err
	}

	return spaceID, nil

}
