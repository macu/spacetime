<template>
<form-layout title="Create tag">

	<form-field title="Name for new tag" required>
		<el-input
			v-model="tag"
			:maxlength="$store.getters.tagMaxLength"
			show-word-limit
			size="large"
			:disabled="posting"
			/>
	</form-field>

	<form-field v-if="permissions.userAllowPinToParent">
		<el-checkbox v-model="pin" :disabled="posting">
			Pin this tag to the top of the parent space
		</el-checkbox>
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
		permissions: {
			type: Object,
			required: true,
		},
	},
	data() {
		return {
			tag: '',
			pin: false,
		};
	},
	computed: {
		createDisabled() {
			return this.posting || !this.tag.trim();
		},
	},
	methods: {
		submit() {
			if (this.createDisabled) {
				return;
			}
			this.$emit('submit', {
				tag: this.tag.trim(),
				pin: this.permissions.userAllowPinToParent && this.pin,
			});
		},
	},
};
</script>
