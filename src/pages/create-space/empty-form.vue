<template>
<form-layout :title="root ? 'Create empty space' : 'Create subspace'">

	<form-field title="Label for new space" required>
		<template #tip>
			<div>A label is unique and permanent identifier for a given space or subspace. It cannot be changed later.</div>
		</template>
		<el-input
			v-model="label"
			:maxlength="$store.getters.labelMaxLength"
			show-word-limit
			size="large"
			:disabled="posting"
			/>
	</form-field>

	<form-field title="Title for new space">
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
		root: {
			type: Boolean,
			default: false,
		},
		posting: {
			type: Boolean,
			default: false,
		},
	},
	data() {
		return {
			label: '',
			title: '',
		};
	},
	computed: {
		createDisabled() {
			return this.posting || !this.label.trim();
		},
	},
	methods: {
		submit() {
			if (this.createDisabled) {
				return;
			}
			this.$emit('submit', {
				label: this.label.trim(),
				title: this.title.trim(),
			});
		},
	},
};
</script>
