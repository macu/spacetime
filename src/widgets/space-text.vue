<template>
<div class="space-text" v-text="textOutput"/>
</template>

<script>
export default {
	props: {
		space: {
			type: Object,
			required: true,
		},
	},
	computed: {
		textOutput() {
			return (this.space.text || '').trim();
		},
	},
};
</script>

<style lang="scss">
@import '@/styles/vars.scss';

.space-text {
	background-color: $text-bg-color;
	border: $text-border;
	color: $text-color;
	border-radius: $border-radius;
	box-shadow: $text-inner-drop-shadow;

	padding: 20px;

	font-size: 1.2em;
	white-space: pre-wrap;
}
</style>
