package treetime

type NodeClass string

const (
	NodeClassTag      NodeClass = "tag"
	NodeClassLang     NodeClass = "lang"
	NodeClassType     NodeClass = "type"
	NodeClassField    NodeClass = "field"
	NodeClassCategory NodeClass = "category"
	NodeClassPost     NodeClass = "post"
	NodeClassComment  NodeClass = "comment"
)
