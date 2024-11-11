package treetime

import (
	"strings"
	"time"
)

type NodePostBlock struct {
	Type      string  `json:"type"`                // "text" or "link"
	Text      *string `json:"text,omitempty"`      // text blocks
	URL       *string `json:"url,omitempty"`       // link blocks
	ImageData *string `json:"imageData,omitempty"` // link blocks
}

type NodeContent struct {
	Title       *string          `json:"title,omitempty"`       // categories, types, fields, tags, langs
	Description *string          `json:"description,omitempty"` // categories, types, and fields
	Text        *string          `json:"text,omitempty"`        // comments
	Blocks      *[]NodePostBlock `json:"blocks,omitempty"`      // posts
	LangNodeID  *uint            `json:"langNodeID"`            // translation
}

type NodeCreator struct {
	ID          uint    `json:"id"`
	DisplayName *string `json:"displayName,omitempty"`
}

type NodeHeader struct {
	ID        uint          `json:"id"`
	Class     string        `json:"class"`
	Content   NodeContent   `json:"content"`
	Path      *[]NodeHeader `json:"path,omitempty"`
	Tags      *[]NodeTag    `json:"tags,omitempty"`
	OwnerType *string       `json:"ownerType,omitempty"`
	Creator   *NodeCreator  `json:"creator,omitempty"`
	CreatedAt *time.Time    `json:"createdAt,omitempty"`
	IsDeleted *bool         `json:"isDeleted,omitempty"`
}

type LangHeader struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type NodeTag struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	UserVote     *string `json:"userVote"`
	SupportRatio float32 `json:"supportRatio"`
}

func (b NodeContent) extractTextForTSVector() string {
	builder := strings.Builder{}

	if b.Title != nil {
		builder.WriteString(*b.Title)
	}

	if b.Description != nil {
		builder.WriteString(" ")
		builder.WriteString(*b.Description)
	}

	if b.Text != nil {
		builder.WriteString(" ")
		builder.WriteString(*b.Text)
	}

	if b.Blocks != nil {
		for _, block := range *b.Blocks {
			if block.Text != nil {
				builder.WriteString(" ")
				builder.WriteString(*block.Text)
			}
		}
	}

	return strings.TrimSpace(builder.String())
}
