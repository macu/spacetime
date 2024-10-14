package treetime

import (
	"fmt"
	"time"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func SetNodeContentVote(conn db.DBConn, auth ajax.Auth, nodeID, contentID uint, vote string) error {

	if !IsValidVote(vote) {
		return fmt.Errorf("invalid vote: %s", vote)
	}

	_, err := conn.Exec(`INSERT INTO tree_node_content_vote
		(node_id, content_id, created_at, created_by, vote)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (created_by, node_id, content_id)
		DO UPDATE SET vote = $5, created_at = $3`,
		nodeID, contentID, time.Now(), auth.UserID, vote,
	)
	if err != nil {
		return fmt.Errorf("inserting tree_node_content_vote: %w", err)
	}

	return nil

}
