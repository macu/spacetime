<template>
<div class="space-title">
	<checkin-button :space="space" @check-in="$emit('check-in')"/>
	<span class="space-title-text" v-text="titleOutput"/>
</div>
</template>

<script>
export default {
	emits: [
		'check-in',
	],
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
.title-space {
	display: flex;
	flex-direction: row;
	align-items: center;
	flex-wrap: nowrap;

	>.space-title-text {
		font-size: 1.5em;
		font-weight: bold;
		white-space: nowrap;
	}
}
</style>
