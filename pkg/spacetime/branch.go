package spacetime

import (
	"database/sql"
	"fmt"

	"spacetime/pkg/utils/ajax"
)

func CreateBranchSpace(tx *sql.Tx, auth ajax.Auth,
	parentID *uint, label string,
) (*Space, error) {

	var space = &Space{
		ParentID:  parentID,
		SpaceType: SpaceTypeBranch,
	}

	uniqueTextID, err := GetOrCreateUniqueTextId(tx, label)
	if err != nil {
		return nil, fmt.Errorf("get or create unique text id: %w", err)
	}

	// Create space
	if err = CreateSpace(tx, auth, space, parentID, SpaceTypeBranch); err != nil {
		return nil, fmt.Errorf("create branch space: %w", err)
	}

	if parentID == nil {
		// 0 representing root for uniqueness check in branch_space table
		parentID = new(uint)
	}

	if _, err = tx.Exec(`INSERT INTO branch_space
			(parent_id, space_id, label_text_id)
			VALUES ($1, $2, $3)`,
		parentID, space.ID, uniqueTextID,
	); err != nil {
		return nil, fmt.Errorf("insert branch space: %w", err)
	}

	return space, nil

}
