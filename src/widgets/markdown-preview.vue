<template>
<div class="markdown-preview">

	<loading v-if="loadingPreview" message="Loading markdown preview..."/>

	<div v-else-if="showingPreview" v-html="preview"/>

	<el-button v-else :disabled="previewDisabled" @click="loadPreview()">
		Preview markdown
	</el-button>

</div>
</template>

<script>
export default {
	props: {
		content: {
			type: String,
			required: true,
		},
	},
	data() {
		return {
			preview: '',
			loadingPreview: false,
		};
	},
	computed: {
		showingPreview() {
			return !!this.preview.trim();
		},
		previewDisabled() {
			return !this.content.trim();
		},
	},
	methods: {
		loadPreview() {
			// TODO
		},
	},
};
</script>

<style lang="scss">
.markdown-preview {
	>div {
		padding: 20px;
		background-color: lightblue;
		border: thin solid blue;
	}
}
</style>
