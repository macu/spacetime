package spacetime

import (
	"database/sql"
	"fmt"
	"time"

	"spacetime/pkg/utils/ajax"
	"spacetime/pkg/utils/db"
)

func CreateSpace(conn db.DBConn, auth ajax.Auth,
	space *Space, parentID uint, spaceType string,
) error {

	err := conn.QueryRow(`INSERT INTO space
			(parent_id, space_type, created_at, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, created_by`,
		parentID, spaceType, time.Now(), auth.UserID,
	).Scan(&space.ID, &space.CreatedAt, &space.CreatedBy)

	if err != nil {
		return fmt.Errorf("insert space: %w", err)
	}

	return nil

}

func CreateSubspace(conn db.DBConn, auth ajax.Auth,
	parentID *uint, label string,
) (*Space, error) {

	var space = &Space{
		ParentID:  parentID,
		SpaceType: SpaceTypeSpace,
	}

	err := conn.QueryRow(`INSERT INTO space
			(parent_id, space_type, created_at, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, created_by`,
		parentID, space.SpaceType, time.Now(), auth.UserID,
	).Scan(&space.ID, &space.CreatedAt, &space.CreatedBy)

	if err != nil {
		return nil, fmt.Errorf("insert space: %w", err)
	}

	uniqueTextID, err := GetOrCreateUniqueTextId(conn, label)
	if err != nil {
		return nil, fmt.Errorf("get or create unique text id: %w", err)
	}

	if parentID == nil {
		parentID = new(uint) // 0 for subspace table
	}

	_, err = conn.Exec(`INSERT INTO subspace
		(parent_id, space_id, label_unique_text_id)
		VALUES ($1, $2, $3)`,
		parentID, space.ID, uniqueTextID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert subspace: %w", err)
	}

	return space, nil

}

func GetSpace(conn db.DBConn, spaceID uint) (*Space, error) {

	var space = &Space{
		ID: spaceID,
	}

	err := conn.QueryRow(`SELECT parent_id, space_type, created_at, created_by
		FROM space WHERE id = $1`,
		spaceID,
	).Scan(&space.ParentID, &space.SpaceType, &space.CreatedAt, &space.CreatedBy)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("select space: %w", err)
	}

	return space, nil

}
