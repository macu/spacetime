<template>
<div class="space-tag flex-row nowrap">
	<checkin-button
		v-if="showCheckin"
		:space="space"
		@check-in="$emit('check-in')"
		size="small"
		/>
	<material-icon :icon="icon"/>
	<span
		@click="$emit('click-tag')"
		class="text"
		v-text="tagOutput"
		/>
</div>
</template>

<script>
import CheckinButton from './checkin-button.vue';

import {
	SPACE_TYPE_ICONS,
} from '@/const.js';

export default {
	emits: [
		'check-in',
		'click-title',
	],
	components: {
		CheckinButton,
	},
	props: {
		space: {
			type: Object,
			required: true,
		},
		showCheckin: {
			type: Boolean,
			default: true,
		},
	},
	computed: {
		icon() {
			return SPACE_TYPE_ICONS[this.space.spaceType];
		},
		tagOutput() {
			return this.space.text || '';
		},
	},
};
</script>

<style lang="scss">
.space-tag {
	padding: 5px 20px;
	background-color: white;
	>.text {
		font-size: 1.5em;
		font-weight: bold;
		white-space: nowrap;
		cursor: pointer;
	}
}
</style>
