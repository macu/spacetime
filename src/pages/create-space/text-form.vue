<template>
<form-layout title="Create text">

	<form-field title="Title">
		<el-input
			v-model="title"
			:maxlength="$store.getters.titleMaxLength"
			show-word-limit
			size="large"
			:disabled="posting"
			/>
	</form-field>

	<form-field title="Text" required>
		<el-input
			type="textarea"
			v-model="text"
			:maxlength="$store.getters.textMaxLength"
			show-word-limit
			:autosize="{minRows: 3}"
			:disabled="posting"
			/>
	</form-field>

	<form-actions>
		<el-button @click="submit()" type="primary" :disabled="createDisabled">
			Create
		</el-button>
	</form-actions>

</form-layout>
</template>

<script>
export default {
	props: {
		posting: {
			type: Boolean,
			default: false,
		},
	},
	data() {
		return {
			title: '',
			text: '',
		};
	},
	computed: {
		createDisabled() {
			return this.posting || !this.text.trim();
		},
	},
	methods: {
		submit() {
			if (this.createDisabled) {
				return;
			}
			this.$emit('submit', {
				title: this.title.trim(),
				text: this.text.trim(),
			});
		},
	},
};
</script>
