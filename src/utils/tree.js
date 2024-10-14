import {
	NODE_CLASS,
	SYSTEM_NODE_KEYS,
} from '@/const.js';

const SYSTEM_KEYS = Object.values(SYSTEM_NODE_KEYS);

export function getPathKeyScope(path) {
	if (!path || !path.length) {
		return null;
	}
	for (let i = 0; i < path.length; i++) {
		let key = path[i].key;
		if (key) {
			if (SYSTEM_KEYS.includes(key)) {
				return key;
			}
		}
	}
	return null;
}

export function getPathClassScope(path) {
	if (!path || !path.length) {
		return null;
	}
	let nodeClass = null;
	for (let i = 0; i < path.length; i++) {
		nodeClass = path[i].class;
		switch (nodeClass) {
			case NODE_CLASS.CATEGORY:
			case NODE_CLASS.POST:
				continue;
			case NODE_CLASS.TAG:
			case NODE_CLASS.TYPE:
			case NODE_CLASS.FIELD:
			case NODE_CLASS.LANG:
			case NODE_CLASS.COMMENT:
				return nodeClass;
		}
	}
	return nodeClass;
}

export function allowCreateLang(path) {
	let keyScope = getPathKeyScope(path);
	let classScope = getPathClassScope(path);
	return keyScope === SYSTEM_NODE_KEYS.LANGS && classScope === NODE_CLASS.CATEGORY;
}

export function allowCreateTag(path) {
	let keyScope = getPathKeyScope(path);
	let classScope = getPathClassScope(path);
	return keyScope === SYSTEM_NODE_KEYS.TAGS && classScope === NODE_CLASS.CATEGORY;
}

export function allowCreateType(path) {
	let keyScope = getPathKeyScope(path);
	let classScope = getPathClassScope(path);
	return keyScope === SYSTEM_NODE_KEYS.TYPES && classScope === NODE_CLASS.CATEGORY;
}

export function allowCreateField(path) {
	let keyScope = getPathKeyScope(path);
	let classScope = getPathClassScope(path);
	return keyScope === SYSTEM_NODE_KEYS.TYPES && classScope === NODE_CLASS.TYPE;
}

export function allowCreateCategory(path) {
	let classScope = getPathClassScope(path);
	return classScope === null || classScope === NODE_CLASS.CATEGORY;
}

export function allowCreatePost(path) {
	let keyScope = getPathKeyScope(path);
	let classScope = getPathClassScope(path);
	return (keyScope === null || keyScope === SYSTEM_NODE_KEYS.TREETIME) &&
		(classScope === null || classScope === NODE_CLASS.CATEGORY);
}

export function allowCreateComment(path) {
	return true; // allow anywhere
}
