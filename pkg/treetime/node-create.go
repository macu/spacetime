package treetime

import (
	"fmt"
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

func VerifyCreateNodeAllowed(db db.DBConn, parentID uint, createClass NodeClass) (bool, error) {

	if parentID == 0 {
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

	parentClass, parentKey, err := LoadNodeMeta(db, parentID)

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
			path, err := LoadNodeParentPath(db, nil, parentID)
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
			path, err := LoadNodeParentPath(db, nil, parentID)
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
			path, err := LoadNodeParentPath(db, nil, parentID)
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
			path, err := LoadNodeParentPath(db, nil, parentID)
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

func CreateNode(db db.DBConn, auth ajax.Auth, parentID uint, class NodeClass, title, content string) (*NodeHeader, error) {

	allowed, err := VerifyCreateNodeAllowed(db, parentID, class)

	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, fmt.Errorf("not allowed")
	}

	return nil, fmt.Errorf("not implemented")

}
