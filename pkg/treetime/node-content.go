package treetime

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

// remove newlines and replace all whitespace with a single space
var sanitizeWhitespaceRegex = regexp.MustCompile(`\s+`)

func FormatTitle(content *string) *string {
	if content != nil {
		newContent := sanitizeWhitespaceRegex.ReplaceAllString(strings.TrimSpace(*content), " ")
		return &newContent
	}
	return content
}

func FormatDescription(content *string) *string {
	if content != nil {
		newContent := strings.TrimSpace(*content)
		return &newContent
	}
	return content
}

func SanitizeNodeContent(class string, content *NodeContent) {
	switch class {

	case NodeClassCategory, NodeClassType, NodeClassField:
		content.Title = FormatTitle(content.Title)
		content.Description = FormatDescription(content.Description)
		content.Text = nil
		content.Blocks = nil

	case NodeClassTag:
		content.Title = FormatTitle(content.Title)
		content.Description = nil
		content.Text = nil
		content.Blocks = nil

	case NodeClassPost:
		content.Title = FormatTitle(content.Title)
		content.Description = nil
		content.Text = nil
		if content.Blocks != nil {
			filteredBlocks := []NodePostBlock{}
			for i := range *content.Blocks {
				block := (*content.Blocks)[i]
				if block.Type == NodePostBlockTypeText {
					block.Text = FormatDescription(block.Text)
					if block.Text != nil && *block.Text != "" {
						filteredBlocks = append(filteredBlocks, block)
					}
				}
			}
			content.Blocks = &filteredBlocks
		}

	case NodeClassComment:
		content.Title = FormatTitle(content.Title)
		content.Description = nil
		content.Text = FormatDescription(content.Text)
		content.Blocks = nil

	}
}

func IsValidNodeContent(class string, content NodeContent) bool {
	switch class {

	case NodeClassCategory, NodeClassType, NodeClassField:
		if content.Title == nil || *content.Title == "" {
			return false
		}
		if content.Text != nil {
			return false
		}
		if content.Blocks != nil {
			return false
		}

	case NodeClassTag:
		if content.Title == nil || *content.Title == "" {
			return false
		}
		if content.Description != nil {
			return false
		}
		if content.Text != nil {
			return false
		}
		if content.Blocks != nil {
			return false
		}

	case NodeClassPost:
		if content.Title == nil || *content.Title == "" {
			return false
		}
		if content.Description != nil {
			return false
		}
		if content.Text != nil {
			return false
		}
		if content.Blocks == nil || len(*content.Blocks) == 0 {
			return false
		}

	case NodeClassComment:
		if content.Text == nil || *content.Text == "" {
			return false
		}
		if content.Description != nil {
			return false
		}
		if content.Blocks != nil {
			return false
		}

	case NodeClassLang:
		if content.Title == nil || *content.Title == "" {
			return false
		}
		if content.Description != nil {
			return false
		}
		if content.Text != nil {
			return false
		}
		if content.Blocks != nil {
			return false
		}

	}

	return true
}

func LoadContentForNodes(conn db.DBConn, auth *ajax.Auth, headers []NodeHeader) error {

	if len(headers) == 0 {
		return nil
	}

	var nodeIDs []uint
	for _, header := range headers {
		nodeIDs = append(nodeIDs, header.ID)
	}

	var args = []interface{}{}

	rows, err := conn.Query(`SELECT
		tree_node_content.node_id, tree_node_content.content_json
		FROM tree_node_content
		WHERE `+db.In("tree_node_content.node_id", &args, nodeIDs),
		args...,
	)

	if err != nil {
		return fmt.Errorf("loading node content: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id uint
		var jsonContent string
		err = rows.Scan(&id, &jsonContent)
		if err != nil {
			return fmt.Errorf("scanning node content: %w", err)
		}

		for i := range headers {
			if headers[i].ID == id {
				err = json.Unmarshal([]byte(jsonContent), &headers[i].Content)
				if err != nil {
					return fmt.Errorf("unmarshalling node content: %w", err)
				}
				break
			}
		}
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
