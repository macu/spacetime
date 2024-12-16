<template>
<span class="space-type" v-text="typeOutput"/>
</template>

<script>
import {
	SPACE_TYPES,
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
	},
};
</script>

<style lang="scss">
.space-type {
	color: darkblue;
	text-shadow: 1px 1px 1px blue;
}
</style>
