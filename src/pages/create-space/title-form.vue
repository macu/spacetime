<template>
<form-layout title="Create title">

	<form-field title="New title" required>
		<el-input
			v-model="title"
			:maxlength="$store.getters.titleMaxLength"
			show-word-limit
			size="large"
			:diaabled="posting"
			/>
	</form-field>

	<form-actions>
		<el-button @click="create()" type="primary" :disabled="createDisabled">
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
