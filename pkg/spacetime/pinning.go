package spacetime

import (
	"database/sql"
	"fmt"
	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func AllowsPinToParent(
	conn *sql.DB, auth ajax.Auth,
	userID, parentID uint) (bool, error) {

	// Load parent path
	path, err := LoadParentPath(conn, &auth, parentID)
	if err != nil {
		return false, fmt.Errorf("get space parent path: %w", err)
	}

	if len(path) == 0 {
		// no pinning at root level
		return false, nil
	}

	rootSpace := path[0]
	if rootSpace.SpaceType == SpaceTypeUser {
		// only user has pinning permission within their own personal space
		return rootSpace.CreatedBy == auth.UserID, nil
	}

	// Check if nearest author context space in path is created by the user
	reversedPath := make([]*Space, len(path))
	for i, j := len(path)-1, 0; i >= 0; i, j = i-1, j+1 {
		reversedPath[j] = path[i]
	}
	for _, node := range reversedPath {
		if node.SpaceType == SpaceTypeText {
			// Allow user to pin items within their own public text spaces,
			// even under another user's text space
			return node.CreatedBy == auth.UserID, nil
		}
	}

	return false, nil

}

func AllowsPinningUnderCreateSpace(
	conn *sql.DB, auth ajax.Auth,
	userID, parentID uint, createSpaceType string) (bool, error) {

	// Load parent path
	path, err := LoadParentPath(conn, &auth, parentID)
	if err != nil {
		return false, fmt.Errorf("get space parent path: %w", err)
	}

	if len(path) == 0 {
		// no pinning at root level
		return false, nil
	}

	rootSpace := path[0]
	if rootSpace.SpaceType == SpaceTypeUser {
		// only user has pinning permission within their own personal space
		return rootSpace.CreatedBy == auth.UserID, nil
	}

	if createSpaceType == SpaceTypeText {
		// Allow user to pin under any public text they create
		return true, nil
	}

	// Check if nearest author context space in path is created by the user
	reversedPath := make([]*Space, len(path))
	for i, j := len(path)-1, 0; i >= 0; i, j = i-1, j+1 {
		reversedPath[j] = path[i]
	}
	for _, node := range reversedPath {
		if node.SpaceType == SpaceTypeText {
			// Allow user to pin items within their own public text spaces,
			// even under another user's text space
			return node.CreatedBy == auth.UserID, nil
		}
	}

	return false, nil

}

func PinSpace(conn db.DBConn, auth ajax.Auth, space *Space) error {

	// Get max existing pinned branch order number
	var maxOrderNumber *int
	err := conn.QueryRow(`SELECT MAX(order_number) FROM user_space_config
			INNER JOIN space ON space.id = user_space_config.space_id
			WHERE space.parent_id = $1`,
		space.ParentID,
	).Scan(&maxOrderNumber)
	if err != nil {
		return fmt.Errorf("get max pinned branch order number: %w", err)
	}

	// Insert new pinned branch with order number after max existing
	_, err = conn.Exec(`INSERT INTO user_space_config
			(space_id, user_id, order_number)
			VALUES ($1, $2, COALESCE($3, -1) + 1)
			ON CONFLICT DO NOTHING`,
		space.ID, auth.UserID, maxOrderNumber,
	)
	if err != nil {
		return fmt.Errorf("insert pinned branch space config: %w", err)
	}

	return nil

}

func PinSpaces(conn db.DBConn, auth ajax.Auth, parentID uint, spaces []*Space) error {

	for _, space := range spaces {
		err := PinSpace(conn, auth, space)
		if err != nil {
			return fmt.Errorf("pin space: %w", err)
		}
	}

	return nil

}

func UnpinSpace(conn db.DBConn, auth *ajax.Auth, spaceID uint) error {

	_, err := conn.Exec(`DELETE FROM user_space_config
		WHERE space_id = $1 AND user_id = $2`,
		spaceID, auth.UserID,
	)
	if err != nil {
		return fmt.Errorf("delete pinned space config: %w", err)
	}

	return nil

}

func ReorderPin(conn *sql.DB, auth *ajax.Auth, pinnedSpace *Space, orderNumber uint) error {

	return db.InTransaction(conn, func(tx *sql.Tx) error {

		// Get existing order number for pinned space
		var existingOrderNumber *uint
		err := tx.QueryRow(`SELECT order_number FROM user_space_config
			WHERE space_id = $1 AND user_id = $2`,
			pinnedSpace.ID, auth.UserID,
		).Scan(&existingOrderNumber)
		if err != nil {
			return fmt.Errorf("get existing pinned space order number: %w", err)
		} else if existingOrderNumber == nil {
			return fmt.Errorf("pinned space config not found")
		}

		if *existingOrderNumber == orderNumber {
			// No change needed
			return nil
		}

		// Shift order numbers down following gap
		_, err = tx.Exec(`UPDATE user_space_config
			SET order_number = order_number - 1
			WHERE user_id = $1 AND space_id IN (
				SELECT space.id FROM space
				INNER JOIN user_space_config ON user_space_config.space_id = space.id
				WHERE space.parent_id = $2 AND user_space_config.order_number > $3
			)`,
			auth.UserID, pinnedSpace.ParentID, *existingOrderNumber,
		)
		if err != nil {
			return fmt.Errorf("shift down pinned branch order numbers: %w", err)
		}

		// Shift order numbers up to make space for moved pinned space
		_, err = tx.Exec(`UPDATE user_space_config
			SET order_number = order_number + 1
			WHERE user_id = $1 AND space_id IN (
				SELECT space.id FROM space
				INNER JOIN user_space_config ON user_space_config.space_id = space.id
				WHERE space.parent_id = $2 AND user_space_config.order_number >= $3
			)`,
			auth.UserID, pinnedSpace.ParentID, orderNumber,
		)
		if err != nil {
			return fmt.Errorf("shift up pinned branch order numbers: %w", err)
		}

		// Update config for moved pinned space with new order number
		_, err = tx.Exec(`UPDATE user_space_config
			SET order_number = $3
			WHERE space_id = $1 AND user_id = $2`,
			pinnedSpace.ID, auth.UserID, orderNumber,
		)
		if err != nil {
			return fmt.Errorf("update moved pinned space config: %w", err)
		}

		// Reorder all sequentially
		_, err = tx.Exec(`WITH ordered_pins AS (
			SELECT space.id, ROW_NUMBER() OVER (ORDER BY user_space_config.order_number) - 1 AS new_order_number
			FROM space
			INNER JOIN user_space_config ON user_space_config.space_id = space.id
			WHERE space.parent_id = $1 AND user_space_config.user_id = $2
		)
		UPDATE user_space_config
		SET order_number = ordered_pins.new_order_number
		FROM ordered_pins
		WHERE user_space_config.space_id = ordered_pins.id AND user_space_config.user_id = $2`,
			pinnedSpace.ParentID, auth.UserID,
		)
		if err != nil {
			return fmt.Errorf("reorder pinned branch order numbers: %w", err)
		}

		return nil

	})

}
