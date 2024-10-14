<template>
<el-input
	:modelValue="modelValue"
	@update:modelValue="emitUpdate"
	type="textarea"
	:placeholder="placeholder"
	:maxlength="maxlength"
	:autosize="{minRows: 1}"
	resize="none"
	clearable
	show-word-limit
	/>
</template>

<script>
const CONDENSE_WHITESPACE = /(?:[\n\r\t\v\f]|\s{2,})/g;

export default {
	props: {
		modelValue: {
			type: String,
			required: true,
		},
		maxlength: {
			type: Number,
			required: false,
			default: null,
		},
		placeholder: {
			type: String,
			required: false,
			default: '',
		},
	},
	watch: {
		modelValue: {
			immediate: true,
			handler(value) {
				if (CONDENSE_WHITESPACE.test(value)) {
					this.emitUpdate(value);
				}
			},
		},
	},
	methods: {
		emitUpdate(value) {
			// Replace all whitespace with a single space
			this.$emit('update:modelValue', value.replace(/\s+/g, ' '));
		},
	},
};
</script>
