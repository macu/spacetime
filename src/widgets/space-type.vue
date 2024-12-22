<template>
<span class="space-type flex-row nowrap">
	<material-icon :icon="icon"/>
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
		type: {
			type: String,
			required: true,
		},
	},
	computed: {
		typeOutput() {
			switch (this.type) {
				case SPACE_TYPES.SPACE:
					return 'Space';
				case SPACE_TYPES.CHECK_IN:
					return 'Check in';
				case SPACE_TYPES.TITLE:
					return 'Title';
				case SPACE_TYPES.TAG:
					return 'Tag';
				case SPACE_TYPES.TEXT:
					return 'Text';
				case SPACE_TYPES.NAKED_TEXT:
					return 'Naked text';
				case SPACE_TYPES.STREAM_OC:
					return 'Stream of consciousness';
				case SPACE_TYPES.JSON_AR:
					return 'JSON attribute';
			}
			return 'Unknown';
		},
		icon() {
			return SPACE_TYPE_ICONS[this.type];
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
