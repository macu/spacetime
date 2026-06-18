package spacetime

const LabelMaxLength = 128
const TextMaxLength = 2048
const TitleMaxLength = 256
const TagMaxLength = 128
const NakedTextMaxDeltasSoft = TextMaxLength * 3
const NakedTextMaxDeltas = TextMaxLength * 4

const SpaceTypeUser = "user"
const SpaceTypeBranch = "branch"
const SpaceTypeText = "text"
const SpaceTypeTag = "tag"
const SpaceTypeLink = "link"

// const SpaceTypeStream = "stream-of-consciousness"
// const SpaceTypeJson = "json-attribute"

func IsValidSpaceType(spaceType string) bool {
	switch spaceType {

	case SpaceTypeBranch,
		SpaceTypeText,
		SpaceTypeLink,
		SpaceTypeTag:
		return true

	default:
		// Not yet implemented
		return false
	}
}

func IsValidRootSpaceType(spaceType string) bool {
	switch spaceType {

	case SpaceTypeBranch:
		return true

	default:
		return false
	}
}

func IsAuthorContextSpace(spaceType string) bool {
	switch spaceType {

	case SpaceTypeUser, SpaceTypeText:
		return true

	default:
		return false
	}
}
