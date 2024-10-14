package treetime

import (
	"treetime/pkg/utils/ajax"
	"treetime/pkg/utils/db"
)

type NodeTag struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	UserVote     *string `json:"userVote"`
	SupportRatio float32 `json:"supportRatio"`
}

type NodeLangTag struct {
	ID           uint     `json:"id"`
	Title        string   `json:"title"`
	UserVote     *string  `json:"userVote"`
	SupportRatio float32  `json:"supportRatio"`
	LangCodes    []string `json:"langCodes"`
}

func LoadNodeTags(db db.DBConn, auth *ajax.Auth, nodeId uint) ([]NodeTag, error) {
	var tags []NodeTag

	return tags, nil
}
