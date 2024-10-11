package treetime

type NodeSearchParams struct {
	LimitToClass NodeClass
	ExcludeClass NodeClass
}

type TagQueryParams struct {
	LangTagIDs    []uint
	PreferTagIDs  []uint // prefer results with these tags, in order
	OrderByTagIDs []uint // order results by these tags, combined
	ExcludeTagIDs []uint // exclude results with these tags
}
