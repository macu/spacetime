<template>
<span class="space-type flex-row nowrap">
	<material-icon :icon="icon"/>
	<div v-if="showPin && isPinned" class="pinned-indicator">
		<el-tooltip content="This space is pinned by the author" placement="top">
			<material-icon icon="keep" class="pinned"/>
		</el-tooltip>
	</div>
	<span v-text="typeOutput"/>
</span>
</template>

<script>
import {
	SPACE_TYPES,
	SPACE_TYPE_ICONS,
} from '@/const.js';

export default {
	props: {
		space: {
			type: Object,
			required: true,
		},
		showPin: {
			type: Boolean,
			default: false,
		},
	},
	computed: {
		typeOutput() {
			switch (this.space.spaceType) {
				case SPACE_TYPES.USER:
					return this.space.authorHandle ||
						this.space.authorDisplayName || 'User';
				case SPACE_TYPES.BRANCH:
					return 'Branch';
				case SPACE_TYPES.TAG:
					return 'Tag';
				case SPACE_TYPES.TEXT:
					return 'Text';
				case SPACE_TYPES.STREAM_OC:
					return 'Stream of consciousness';
				case SPACE_TYPES.JSON_AR:
					return 'JSON attribute';
			}
			return 'Unknown';
		},
		icon() {
			return SPACE_TYPE_ICONS[this.space.spaceType];
		},
		isPinned() {
			return this.space.isPinned;
		},
	},
};
</script>

<style lang="scss">
.space-type {
	color: darkblue;
	text-shadow: 2px 2px 3px white;
	font-size: larger;
}
</style>
