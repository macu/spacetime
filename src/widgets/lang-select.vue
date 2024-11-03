<template>
<el-select
	:modelValue="modelValue"
	@update:modelValue="val => $emit('update:modelValue', val)">
	<el-option
		v-for="lang in langs"
		:key="lang.id"
		:label="lang.title"
		:value="lang.id"
		/>
</el-select>
</template>

<script>
export default {
	emits: ['update:modelValue'],
	props: {
		modelValue: {
			type: Number,
			required: false,
		},
	},
	computed: {
		langs() {
			return this.$store.state.langs;
		},
	},
};
</script>
