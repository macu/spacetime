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

<style lang="css">
.space-text {
	background-color: white;
	padding: 20px;

	font-size: 1.2em;
	white-space: pre-wrap;
}
</style>
