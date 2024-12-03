package spacetime

const TitleMaxLength = 128
const TagMaxLength = 128
const TextMaxLength = 1024

const SpaceTypeSpace = "space"
const SpaceTypeCheckin = "checkin"
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

	case SpaceTypeSpace, SpaceTypeCheckin, SpaceTypeTitle, SpaceTypeTag,
		SpaceTypeText, SpaceTypeNaked, SpaceTypeStream,
		SpaceTypeJson:
		return true

	case SpaceTypePicture, SpaceTypeAudio, SpaceTypeVideo:
		// Not yet inplemented
		return false

	default:
		return false
	}
}
