package treetime

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func CheckCreateNodeAllowed(db db.DBConn, parentID *uint, createClass string) (bool, error) {

	if parentID == nil {
		switch createClass {

		case NodeClassCategory:
			// Only allow creating categories at the root
			return true, nil

		default:
			return false, nil
		}
	}

	if createClass == NodeClassComment {
		// Comments can be created on any node
		return true, nil
	}

	parentClass, parentKey, err := LoadNodeMeta(db, *parentID)

	if err != nil {
		return false, err
	}

	if parentKey != nil {
		// Check if the parent is a special node
		// that only allows certain types
		switch *parentKey {

		case NodeKeyLangs:
			switch createClass {
			case NodeClassLang, NodeClassCategory:
				return true, nil
			default:
				return false, nil
			}

		case NodeKeyTags:
			switch createClass {
			case NodeClassTag, NodeClassCategory:
				return true, nil
			default:
				return false, nil
			}

		case NodeKeyTypes:
			switch createClass {
			case NodeClassType, NodeClassCategory:
				return true, nil
			default:
				return false, nil
			}
		}
	}

	switch parentClass {

	case NodeClassCategory:
		switch createClass {

		case NodeClassCategory:
			// Always allow category within category
			return true, nil

		case NodeClassPost:
			// Check all parents are categories
			path, err := LoadNodeParentPath(db, nil, *parentID, false)
			if err != nil {
				return false, err
			}
			for _, header := range path {
				if header.Class != NodeClassCategory {
					return false, nil
				}
			}
			return true, nil

		case NodeClassType:
			// Check parent is "types"
			path, err := LoadNodeParentPath(db, nil, *parentID, false)
			if err != nil {
				return false, err
			}
			for _, header := range path {
				if header.Key != nil && *header.Key == NodeKeyTypes {
					return true, nil
				} else if header.Class != NodeClassCategory {
					break
				}
			}
			return false, nil

		case NodeClassTag:
			// Check parent is "tags"
			path, err := LoadNodeParentPath(db, nil, *parentID, false)
			if err != nil {
				return false, err
			}
			for _, header := range path {
				if header.Key != nil && *header.Key == NodeKeyTags {
					return true, nil
				} else if header.Class != string(NodeClassCategory) {
					break
				}
			}
			return false, nil

		case NodeClassLang:
			// Check parent is "langs"
			path, err := LoadNodeParentPath(db, nil, *parentID, false)
			if err != nil {
				return false, err
			}
			for _, header := range path {
				if header.Key != nil && *header.Key == NodeKeyLangs {
					return true, nil
				} else if header.Class != string(NodeClassCategory) {
					break
				}
			}
			return false, nil

		default:
			return false, nil
		}

	case NodeClassType:
		switch createClass {
		case NodeClassField:
			return true, nil
		default:
			return false, nil
		}

	default:
		return false, nil
	}

}

func CreateNode(conn *sql.DB, auth ajax.Auth, parentID *uint, class, title, body string) (*NodeHeader, error) {

	if !IsValidNodeClass(class) {
		return nil, fmt.Errorf("invalid node class: %s", class)
	}

	title = FormatTitle(title)
	if !CheckContentLength(class, ContentTypeTitle, title) {
		return nil, fmt.Errorf("invalid title length for class: %s", class)
	}

	body = strings.TrimSpace(body)
	if !CheckContentLength(class, ContentTypeBody, body) {
		return nil, fmt.Errorf("invalid body length for class: %s", class)
	}

	allowed, err := CheckCreateNodeAllowed(conn, parentID, class)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("create node location not allowed: %s under node %d", class, *parentID)
	}

	var node = NodeHeader{}
	err = db.InTransaction(conn, func(tx *sql.Tx) error {

		err = tx.QueryRow(`INSERT INTO tree_node
			(node_class, parent_id, created_at, created_by)
			VALUES ($1, $2, $3, $4)
			RETURNING id, node_class`,
			class, parentID, time.Now(), auth.UserID,
		).Scan(&node.ID, &node.Class)
		if err != nil {
			return fmt.Errorf("error creating node: %w", err)
		}

		err = SetNodeVote(tx, auth, parentID, node.ID, VoteAgree)
		if err != nil {
			return fmt.Errorf("error setting node vote: %w", err)
		}

		if title != "" {
			var titleContentID uint
			err = tx.QueryRow(`INSERT INTO tree_node_content
				(node_id, content_type, text_content, created_at, created_by)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, text_content`,
				node.ID, ContentTypeTitle, title, time.Now(), auth.UserID,
			).Scan(&titleContentID, &node.Title)
			if err != nil {
				return fmt.Errorf("error creating node content title: %w", err)
			}

			err = SetNodeContentVote(tx, auth, node.ID, titleContentID, VoteAgree)
			if err != nil {
				return fmt.Errorf("error setting node content vote: %w", err)
			}
		}

		if body != "" {
			var bodyContentID uint
			err = tx.QueryRow(`INSERT INTO tree_node_content
				(node_id, content_type, text_content, created_at, created_by)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, text_content`,
				node.ID, ContentTypeBody, body, time.Now(), auth.UserID,
			).Scan(&bodyContentID, &node.Body)
			if err != nil {
				return fmt.Errorf("error creating node content body: %w", err)
			}

			err = SetNodeContentVote(tx, auth, node.ID, bodyContentID, VoteAgree)
			if err != nil {
				return fmt.Errorf("error setting node content vote: %w", err)
			}
		}

		return nil

	})
	if err != nil {
		return nil, err
	}

	return &node, nil

}
