package treetime

import (
	"encoding/json"
	"fmt"

	"treetime/pkg/utils/db"
)

func LoadLangs(conn db.DBConn) ([]LangHeader, error) {

	var langs = make([]LangHeader, 0)

	// TODO Order by most voted, with user preferred langs at top
	rows, err := conn.Query(`SELECT tree_node.id, tree_node_content.content_json
		FROM tree_node
		INNER JOIN
			tree_node_content
		ON
			tree_node.id = tree_node_content.node_id
			AND tree_node_content.lang_node_id = tree_node.id
		WHERE
			tree_node.node_class = $1`,
		NodeClassLang,
	)

	if err != nil {
		return nil, fmt.Errorf("loading langs: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var lang = LangHeader{}
		var contentJSON string
		var content NodeContent
		err = rows.Scan(&lang.ID, &contentJSON)
		if err != nil {
			return nil, fmt.Errorf("scanning lang: %w", err)
		}
		err = json.Unmarshal([]byte(contentJSON), &content)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling lang: %w", err)
		}
		if content.Title != nil {
			lang.Title = *content.Title
			langs = append(langs, lang)
		}
	}

	return langs, nil

}

func IsValidLangNodeID(conn db.DBConn, id uint) (bool, error) {

	var exists bool
	err := conn.QueryRow(`SELECT EXISTS(SELECT 1 FROM tree_node WHERE id = $1 AND node_class = $2)`,
		id, NodeClassLang,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("checking lang node exists: %w", err)
	}

	return exists, nil

}
