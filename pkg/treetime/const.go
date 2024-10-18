package treetime

const (
	TreeMaxDepth = 50

	CategoryTitleMaxLength       = 100
	CategoryDescriptionMaxLength = 300
	LangTitleMaxLength           = 50
	TagTitleMaxLength            = 50
	TypeTitleMaxLength           = 50
	TypeDescriptionMaxLength     = 200
	PostTitleMaxLength           = 100
	PostBlockMaxLength           = 1024
	CommentMaxLength             = 1024
)

const (
	NodeClassTag      string = "tag"
	NodeClassLang     string = "lang"
	NodeClassType     string = "type"
	NodeClassField    string = "field"
	NodeClassCategory string = "category"
	NodeClassPost     string = "post"
	NodeClassComment  string = "comment"
)

const (
	NodeKeyTreeTime string = "treetime"
	NodeKeyLangs    string = "langs"
	NodeKeyTags     string = "tags"
	NodeKeyTypes    string = "types"
)

const (
	ContentTypeTitle string = "title"
	ContentTypeBody  string = "description"
)

const (
	VoteAgree    string = "agree"
	VoteDisagree string = "disagree"
)

func IsValidNodeClass(class string) bool {
	switch class {
	case NodeClassTag, NodeClassLang, NodeClassType, NodeClassField,
		NodeClassCategory, NodeClassPost, NodeClassComment:
		return true
	}
	return false
}

func IsValidVote(vote string) bool {
	switch vote {
	case VoteAgree, VoteDisagree:
		return true
	}
	return false
}

func IsValidNodeContentType(contentType string) bool {
	switch contentType {
	case ContentTypeTitle, ContentTypeBody:
		return true
	}
	return false
}
