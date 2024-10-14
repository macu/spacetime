package treetime

import (
	"fmt"
	"regexp"
	"strings"

	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

// remove newlines and replace all whitespace with a single space
var titleReplaceWhitespaceRegex = regexp.MustCompile(`(?:[\n\r\t\v\f]|\s{2,})`)

// process whitespace to allow 2 newlines, no tabs, and no leading/trailing whitespace
var bodyExcludedWhitespaceRegex = regexp.MustCompile(`[\r\v\f]`)
var bodyCondenseNewlineRegex = regexp.MustCompile(`(?:\s*[\n]\s*){2,}`)
var bodyReplaceWhitespaceRegex = regexp.MustCompile(`(?:\t|[ ]{2,})`)
var bodyStripHangingWhitespaceRegex = regexp.MustCompile(`(?:\n[ ]+|[ ]+\n)`)

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
	return titleReplaceWhitespaceRegex.ReplaceAllString(content, " ")
}

func FormatBody(content string) string {
	content = strings.TrimSpace(content)
	content = bodyExcludedWhitespaceRegex.ReplaceAllString(content, "")
	content = bodyCondenseNewlineRegex.ReplaceAllString(content, "\n\n")
	content = bodyReplaceWhitespaceRegex.ReplaceAllString(content, " ")
	return bodyStripHangingWhitespaceRegex.ReplaceAllString(content, "\n")
}

func LoadContentTypeForNodes(conn db.DBConn, auth *ajax.Auth, contentType string, headers []NodeHeader) error {

	if !IsValidNodeContentType(contentType) {
		return fmt.Errorf("invalid content type: %s", contentType)
	}

	if len(headers) == 0 {
		return nil
	}

	var nodeIDs []uint
	for _, header := range headers {
		nodeIDs = append(nodeIDs, header.ID)
	}

	var args = []interface{}{contentType}

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
		var content string
		err = rows.Scan(&id, &content)
		if err != nil {
			return fmt.Errorf("scanning node title: %w", err)
		}

		for i := range headers {
			if headers[i].ID == id {
				switch contentType {
				case ContentTypeTitle:
					headers[i].Title = content
				case ContentTypeBody:
					headers[i].Body = content
				}
				break
			}
		}
	}

	return nil

}

func LoadContentForNodes(conn db.DBConn, auth *ajax.Auth, headers []NodeHeader) error {

	err := LoadContentTypeForNodes(conn, auth, ContentTypeTitle, headers)
	if err != nil {
		return err
	}

	err = LoadContentTypeForNodes(conn, auth, ContentTypeBody, headers)
	if err != nil {
		return err
	}

	return nil

}

func LoadContentForNode(conn db.DBConn, auth *ajax.Auth, header *NodeHeader) error {

	var headers = []NodeHeader{*header}

	err := LoadContentForNodes(conn, auth, headers)
	if err != nil {
		return err
	}

	*header = headers[0]

	return nil

}
