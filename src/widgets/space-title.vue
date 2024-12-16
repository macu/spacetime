<template>
<div class="space-title">
	<checkin-button :space="space" @check-in="$emit('check-in')" size="small"/>
	<span class="text" v-text="titleOutput"/>
</div>
</template>

<script>
import CheckinButton from './checkin-button.vue';

export default {
	emits: [
		'check-in',
	],
	components: {
		CheckinButton,
	},
	props: {
		space: {
			type: Object,
			required: true,
		},
	},
	computed: {
		titleOutput() {
			return this.space.text || '';
		},
	},
};
</script>

<style lang="scss">
.space-title {
	display: flex;
	flex-direction: row;
	align-items: center;
	flex-wrap: nowrap;
	column-gap: 10px;

	>.text {
		font-size: 1.5em;
		font-weight: bold;
		white-space: nowrap;
	}
}
</style>
