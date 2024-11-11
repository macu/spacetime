package treetime

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func CreateNode(conn *sql.DB, auth ajax.Auth, parentID, langNodeID *uint, class, ownerType string, content NodeContent) (*NodeHeader, error) {

	if !IsValidNodeClass(class) {
		return nil, fmt.Errorf("invalid node class: %s", class)
	}

	if !IsValidOwnerType(ownerType) {
		return nil, fmt.Errorf("invalid node owner type: %s", ownerType)
	}

	if !IsValidOwnerTypeForClass(class, ownerType) {
		return nil, fmt.Errorf("invalid owner type for class: %s for %s", ownerType, class)
	}

	if !IsValidNodeContent(class, content) {
		return nil, fmt.Errorf("invalid node content for class: %s", class)
	}

	if langNodeID != nil {
		exists, err := IsValidLangNodeID(conn, *langNodeID)
		if err != nil {
			return nil, fmt.Errorf("verifying lang node ID: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("invalid lang node ID: %v", *langNodeID)
		}
	}

	var node = NodeHeader{
		Class:     class,
		Content:   content,
		OwnerType: &ownerType,
	}

	var err error

	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		err = tx.QueryRow(`INSERT INTO tree_node
			(node_class, parent_id, owner_type, created_at, created_by)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`,
			class, parentID, ownerType, time.Now(), auth.UserID,
		).Scan(&node.ID)
		if err != nil {
			return fmt.Errorf("error creating node: %w", err)
		}

		err = SetNodeVote(tx, auth, parentID, node.ID, VoteAgree)
		if err != nil {
			return fmt.Errorf("error setting node vote: %w", err)
		}

		var contentJSON []byte
		contentJSON, err = json.Marshal(content)
		if err != nil {
			return fmt.Errorf("error marshalling node content: %w", err)
		}

		var contentID uint
		err = tx.QueryRow(`INSERT INTO tree_node_content
			(node_id, lang_node_id, content_json, text_search, created_at, created_by)
			VALUES ($1, $2, $3, to_tsvector('pg_catalog.simple', $4), $5, $6)
			RETURNING id`,
			node.ID, langNodeID, contentJSON,
			content.extractTextForTSVector(), time.Now(), auth.UserID,
		).Scan(&contentID)
		if err != nil {
			return fmt.Errorf("error creating node content title: %w", err)
		}

		err = SetNodeContentVote(tx, auth, node.ID, contentID, VoteAgree)
		if err != nil {
			return fmt.Errorf("error setting node content vote: %w", err)
		}

		return nil

	})
	if err != nil {
		return nil, err
	}

	return &node, nil

}
