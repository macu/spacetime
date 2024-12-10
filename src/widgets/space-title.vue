<template>
<div class="title-space">
	{{titleOutput}}
</div>
</template>

<script>
export default {
	props: {
		titleSpace: {
			type: Object,
			required: true,
		},
	},
	computed: {
		titleOutput() {
			return this.titleSpace.text || '';
		},
	},
};
</script>

<style lang="scss">
.title-space {
	font-size: 1.5em;
	font-weight: bold;
}
</style>
