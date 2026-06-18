export const SPACE_TYPES = {
	USER: 'user',
	BRANCH: 'branch',
	TEXT: 'text',
	LINK: 'link',
	TAG: 'tag',
	STREAM_OC: 'stream-of-consciousness',
	JSON_AR: 'json-attribute',
};

export const SPACE_TYPE_ICONS = {
	[SPACE_TYPES.USER]: 'person',
	[SPACE_TYPES.BRANCH]: 'fork_right',
	[SPACE_TYPES.TEXT]: 'description',
	[SPACE_TYPES.LINK]: 'link',
	[SPACE_TYPES.TAG]: 'label',
	[SPACE_TYPES.STREAM_OC]: 'stream',
	[SPACE_TYPES.JSON_AR]: 'code',
};

export const USER_CONTEXT_TYPES = [
	SPACE_TYPES.USER,
	SPACE_TYPES.TEXT,
];
