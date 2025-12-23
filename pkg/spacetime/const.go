package spacetime

const LabelMaxLength = 128
const TitleMaxLength = 256
const TagMaxLength = 128
const TextMaxLength = 2048
const NakedTextMaxDeltasSoft = TextMaxLength * 5
const NakedTextMaxDeltas = TextMaxLength*6 + 1 // may paste textarea maxlength at end

const SpaceTypeSpace = "space"
const SpaceTypeUser = "user"
const SpaceTypeLink = "space-link"
const SpaceTypeCheckin = "check-in"
const SpaceTypeTitle = "title"
const SpaceTypeTag = "tag"
const SpaceTypeText = "text"
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
