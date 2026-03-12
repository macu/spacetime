<template>
<form-layout :title="root ? 'Create space' : 'Create branch'">

	<form-field :title="root ? 'Label for new space' : 'Label for new branch'" required>
		<template #tip>
			<div>A label is unique and permanent identifier for a given space. It cannot be changed later.</div>
		</template>
		<el-input
			v-model="label"
			:maxlength="$store.getters.labelMaxLength"
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
			});
		},
	},
};
</script>
