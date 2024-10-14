package treetime

import (
	"fmt"
	"regexp"
	"strings"

	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

var multipleWhitespaceRegex = regexp.MustCompile(`\s{2,}`)

func CheckContentLength(class, contentType, content string) bool {
	switch class {

	case NodeClassCategory:
		switch contentType {
		case ContentTypeTitle:
			// Title required
			return len(content) > 0 &&
				len(content) <= CategoryTitleMaxLength
		case ContentTypeBody:
			// Body not required
			return len(content) <= CategoryBodyMaxLength
		}

	case NodeClassLang, NodeClassTag:
		switch contentType {
		case ContentTypeTitle:
			// Title required
			return len(content) > 0 &&
				len(content) <= TagTitleMaxLength
		case ContentTypeBody:
			// Body not allowed
			return len(content) == 0
		}

	case NodeClassType, NodeClassField:
		switch contentType {
		case ContentTypeTitle:
			// Title required
			return len(content) > 0 &&
				len(content) <= TagTitleMaxLength
		case ContentTypeBody:
			// Body not required
			return len(content) <= CategoryBodyMaxLength
		}

	case NodeClassPost:
		switch contentType {
		case ContentTypeTitle:
			// Title not required
			return len(content) <= PostTitleMaxLength
		case ContentTypeBody:
			// Body required
			return len(content) > 0 &&
				len(content) <= PostBodyMaxLength
		}

	case NodeClassComment:
		switch contentType {
		case ContentTypeTitle:
			// Title not required
			return len(content) <= PostTitleMaxLength
		case ContentTypeBody:
			// Body required
			return len(content) > 0 &&
				len(content) <= PostBodyMaxLength
		}

	}
	return false
}

func FormatTitle(content string) string {
	content = strings.TrimSpace(content)
	return multipleWhitespaceRegex.ReplaceAllString(content, " ")
}

func LoadNodeTitles(conn db.DBConn, auth *ajax.Auth, headers []NodeHeader) error {

	if len(headers) == 0 {
		return nil
	}

	var nodeIDs []uint
	for _, header := range headers {
		nodeIDs = append(nodeIDs, header.ID)
	}

	var args = []interface{}{ContentTypeTitle}

	rows, err := conn.Query(`SELECT
		tree_node_content.node_id, tree_node_content.text_content
		FROM tree_node_content
		WHERE `+db.In("tree_node_content.node_id", &args, nodeIDs)+`
		AND tree_node_content.content_type = $1`,
		args...,
	)

	if err != nil {
		return fmt.Errorf("loading node titles: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id uint
		var title string
		err = rows.Scan(&id, &title)
		if err != nil {
			return fmt.Errorf("scanning node title: %w", err)
		}

		for i := range headers {
			if headers[i].ID == id {
				headers[i].Title = title
				break
			}
		}
	}

	return nil

}

func LoadNodeTitle(conn db.DBConn, auth *ajax.Auth, header *NodeHeader) error {

	var headers = []NodeHeader{*header}

	err := LoadNodeTitles(conn, auth, headers)
	if err != nil {
		return err
	}

	*header = headers[0]

	return nil

}
