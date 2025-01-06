package spacetime

const TitleMaxLength = 64
const TagMaxLength = 64
const TextMaxLength = 1024

const SpaceTypeSpace = "space"
const SpaceTypeUser = "user"
const SpaceTypeLink = "space-link"
const SpaceTypeCheckin = "check-in"
const SpaceTypeTitle = "title"
const SpaceTypeTag = "tag"
const SpaceTypeText = "text"
const SpaceTypeNaked = "naked-text"
const SpaceTypePicture = "picture"
const SpaceTypeAudio = "audio"
const SpaceTypeVideo = "video"
const SpaceTypeStream = "stream-of-consciousness"
const SpaceTypeJson = "json-attribute"

func IsValidTitle(title string) bool {
	return len(title) > 0 && len(title) <= TitleMaxLength
}

func IsValidTag(tag string) bool {
	return len(tag) > 0 && len(tag) <= TagMaxLength
}

func IsValidText(text string) bool {
	return len(text) > 0 && len(text) <= TextMaxLength
}

func IsValidSpaceType(spaceType string) bool {
	switch spaceType {

	case SpaceTypeSpace,
		SpaceTypeLink,
		SpaceTypeCheckin,
		SpaceTypeTitle, SpaceTypeTag,
		SpaceTypeText:
		return true

	default:
		// Not yet inplemented
		return false
	}
}
