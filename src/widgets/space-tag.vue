<template>
<div class="space-tag flex-row nowrap">
	<checkin-button
		v-if="showCheckin"
		:space="space"
		@check-in="$emit('check-in')"
		size="small"
		/>
	<material-icon icon="label"/>
	<span
		@click="$emit('click-tag')"
		class="text"
		v-text="tagOutput"
		/>
</div>
</template>

<script>
import CheckinButton from './checkin-button.vue';

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
		tagOutput() {
			return this.space.text || '';
		},
	},
};
</script>

<style lang="scss">
.space-tag {
	padding: 5px 10px;
	background-color: white;
	>.text {
		font-size: 1.5em;
		font-weight: bold;
		white-space: nowrap;
		cursor: pointer;
	}
}
</style>
