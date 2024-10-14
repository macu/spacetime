package treetime

import (
	"database/sql"
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

func LoadNodeTitle(conn db.DBConn, auth *ajax.Auth, id uint) (string, error) {

	var title string

	err := conn.QueryRow(`SELECT tree_node_content.text_content
		FROM tree_node_content
		WHERE tree_node_content.node_id = $1
		AND tree_node_content.content_type = $2
		LIMIT 1`,
		id,
		ContentTypeTitle,
	).Scan(&title)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("loading node title: %w", err)
	}

	return title, nil

}
