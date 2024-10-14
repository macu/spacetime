package treetime

const TreeMaxDepth = 50

const CategoryTitleMaxLength = 100
const CategoryDescriptionMaxLength = 500

const (
	NodeClassTag      string = "tag"
	NodeClassLang     string = "lang"
	NodeClassType     string = "type"
	NodeClassField    string = "field"
	NodeClassCategory string = "category"
	NodeClassPost     string = "post"
	NodeClassComment  string = "comment"
)

func IsValidNodeClass(class string) bool {
	switch class {
	case NodeClassTag, NodeClassLang, NodeClassType, NodeClassField,
		NodeClassCategory, NodeClassPost, NodeClassComment:
		return true
	}
	return false
}

const (
	NodeKeyTreeTime string = "treetime"
	NodeKeyLangs    string = "langs"
	NodeKeyTags     string = "tags"
	NodeKeyTypes    string = "types"
)

const (
	ContentTypeTitle string = "title"
	ContentTypeBody  string = "body"
)

func CheckContentLengthLimit(class, contentType, content string) bool {
	switch class {
	case NodeClassCategory:
		switch contentType {
		case ContentTypeTitle:
			return len(content) <= CategoryTitleMaxLength
		case ContentTypeBody:
			return len(content) <= CategoryDescriptionMaxLength
		}
	}
	return false
}
