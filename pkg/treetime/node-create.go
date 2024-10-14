package treetime

import (
	"fmt"
	"strings"
	"time"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func VerifyCreateNodeAllowed(db db.DBConn, parentID *uint, createClass string) (bool, error) {

	if parentID == nil {
		// Only allow creating categories at the root
		switch createClass {

		case NodeClassCategory:
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
			case NodeClassLang:
			case NodeClassCategory:
				return true, nil
			default:
				return false, nil
			}

		case NodeKeyTags:
			switch createClass {
			case NodeClassTag:
			case NodeClassCategory:
				return true, nil
			default:
				return false, nil
			}

		case NodeKeyTypes:
			switch createClass {
			case NodeClassType:
			case NodeClassCategory:
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
			path, err := LoadNodeParentPath(db, nil, *parentID)
			if err != nil {
				return false, err
			}
			for _, header := range path {
				if header.Class != string(NodeClassCategory) {
					return false, nil
				}
			}
			return true, nil

		case NodeClassType:
			// Check parent is "types"
			path, err := LoadNodeParentPath(db, nil, *parentID)
			if err != nil {
				return false, err
			}
			for _, header := range path {
				if header.Key != nil && *header.Key == NodeKeyTypes {
					return true, nil
				} else if header.Class != string(NodeClassCategory) {
					break
				}
			}
			return false, nil

		case NodeClassTag:
			// Check parent is "tags"
			path, err := LoadNodeParentPath(db, nil, *parentID)
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
			path, err := LoadNodeParentPath(db, nil, *parentID)
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

func CreateNode(conn db.DBConn, auth ajax.Auth, parentID *uint, class, title, description string) (*NodeHeader, error) {

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	allowed, err := VerifyCreateNodeAllowed(conn, parentID, class)

	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, fmt.Errorf("create node location not allowed")
	}

	var node = NodeHeader{}
	err = conn.QueryRow(`INSERT INTO tree_node
		(node_class, parent_id, created_at, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, node_class`,
		class, parentID, time.Now(), auth.UserID,
	).Scan(&node.ID, &node.Class)

	if err != nil {
		return nil, fmt.Errorf("error creating node: %w", err)
	}

	err = conn.QueryRow(`INSERT INTO tree_node_content
		(node_id, content_type, text_content, created_at, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING text_content`,
		node.ID, ContentTypeTitle, title, time.Now(), auth.UserID,
	).Scan(&node.Title)

	if err != nil {
		return nil, fmt.Errorf("error creating node content: %w", err)
	}

	if description != "" {
		err = conn.QueryRow(`INSERT INTO tree_node_content
			(node_id, content_type, text_content, created_at, created_by)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING text_content`,
			node.ID, ContentTypeBody, description, time.Now(), auth.UserID,
		).Scan(&node.Description)

		if err != nil {
			return nil, fmt.Errorf("error creating node content: %w", err)
		}
	}

	return &node, nil

}
