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
import {getStorage, setStorage} from '@/utils/storage.js';

const STORAGE_KEY = 'last-lang-select';

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
	watch: {
		modelValue(value) {
			setStorage(STORAGE_KEY, value);
		},
	},
	mounted() {
		if (!this.modelValue) {
			const lastLang = getStorage(STORAGE_KEY);
			if (lastLang) {
				this.$emit('update:modelValue', lastLang);
			}
		}
	},
};
</script>
