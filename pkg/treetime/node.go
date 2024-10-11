package treetime

import (
	"fmt"
	"treetime/pkg/utils/db"
)

type NodeClass string

const (
	NodeClassTag      NodeClass = "tag"
	NodeClassLang     NodeClass = "lang"
	NodeClassType     NodeClass = "type"
	NodeClassField    NodeClass = "field"
	NodeClassCategory NodeClass = "category"
	NodeClassPost     NodeClass = "post"
	NodeClassComment  NodeClass = "comment"
)

func (class NodeClass) String() string {
	return string(class)
}

const NodeKeyTreeTime = "treetime"
const NodeKeyLangs = "langs"
const NodeKeyTags = "tags"
const NodeKeyTypes = "types"

func StringToClass(s string) (NodeClass, error) {
	switch s {
	case "tag":
		return NodeClassTag, nil
	case "lang":
		return NodeClassLang, nil
	case "type":
		return NodeClassType, nil
	case "field":
		return NodeClassField, nil
	case "category":
		return NodeClassCategory, nil
	case "post":
		return NodeClassPost, nil
	case "comment":
		return NodeClassComment, nil
	}
	return "", fmt.Errorf("unknown node class: %s", s)
}

func LoadNodeMeta(db db.DBConn, id uint) (NodeClass, *string, error) {

	var class string
	var key *string

	err := db.QueryRow(`SELECT tree_node.node_class, tree_node_internal_key.internal_key
		FROM tree_node
		LEFt JOIN tree_node_internal_key ON tree_node.id = tree_node_internal_key.node_id
		WHERE tree_node.id = $1`,
		id,
	).Scan(&class, &key)

	if err != nil {
		return "", nil, fmt.Errorf("loading node class: %w", err)
	}

	nodeClass, err := StringToClass(class)

	if err != nil {
		return "", nil, err
	}

	return nodeClass, key, nil

}
