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
	PostBlockCount               = 10
	PostURLMaxLength             = 200
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
	NodePostBlockTypeText string = "text"
)

const (
	OwnerTypeAdmin  string = "admin"
	OwnerTypePublic string = "public"
	OwnerTypeUser   string = "user"
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

func IsValidOwnerType(ownerType string) bool {
	switch ownerType {
	case OwnerTypeAdmin, OwnerTypePublic, OwnerTypeUser:
		return true
	}
	return false
}
