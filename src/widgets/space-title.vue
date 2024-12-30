<template>
<div class="space-title flex-row nowrap" :class="{'ellipsis': ellipsis}">
	<material-icon :icon="icon"/>
	<checkin-button
		v-if="showCheckin"
		:space="space"
		@check-in="$emit('check-in')"
		size="small"
		/>
	<span
		@click="$emit('click-title')"
		class="text"
		v-text="titleOutput"
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
		ellipsis: {
			type: Boolean,
			default: false,
		},
	},
	computed: {
		icon() {
			return SPACE_TYPE_ICONS[this.space.spaceType];
		},
		titleOutput() {
			return this.space.text || '';
		},
	},
};
</script>

<style lang="scss">
@import '@/styles/vars.scss';

.space-title {
	padding: 5px 10px;
	border-radius: $border-radius;
	border: thin solid gray;
	background-color: $title-bg-color;
	color: $title-color;
	border: $title-border;
	overflow: hidden;
	cursor: pointer;

	>.text {
		font-size: 150%;
	}

	&.ellipsis {
		>.text {
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}
	}
}
</style>
