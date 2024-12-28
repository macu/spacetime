<template>
<form-layout title="Create empty space">

	<form-field title="Title for new space" required>
		<el-input
			v-model="title"
			:maxlength="$store.getters.titleMaxLength"
			show-word-limit
			size="large"
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
	emits: ['submit'],
	props: {
		posting: {
			type: Boolean,
			default: false,
		},
	},
	data() {
		return {
			title: '',
		};
	},
	computed: {
		createDisabled() {
			return this.posting || !this.title.trim();
		},
	},
	methods: {
		submit() {
			if (this.createDisabled) {
				return;
			}
			this.$emit('submit', {
				title: this.title.trim(),
			});
		},
	},
};
</script>
