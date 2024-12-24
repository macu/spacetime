<template>
<div class="add-space-tag-widget">
	<div v-if="adding" class="flex-row nowrap">
		<el-input
			v-model="tag"
			:maxlength="$store.getters.tagMaxLength"
			show-word-limit
			size="large">
			<template #prepend>
				Add tag
			</template>
		</el-input>
		<el-button
			@click="addTag()"
			:disabled="addTagDisabled"
			size="large" type="primary">
			Add
		</el-button>
		<el-button
			@click="adding = false"
			size="large" type="warning">
			Cancel
		</el-button>
	</div>
	<el-button v-else @click="adding = true" type="primary" :disabled="disabled">
		Add tag
	</el-button>
</div>
</template>

<script>
import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	emits: [
		'added',
		'update:adding',
	],
	props: {
		parentId: {
			// Tag added under parent space
			type: Number,
			required: true,
		},
	},
	data() {
		return {
			adding: false,
			tag: '',
		};
	},
	computed: {
		addTagDisabled() {
			return !this.tag.trim();
		},
		disabled() {
			return this.$store.getters.createDisabled;
		},
	},
	watch: {
		adding(value) {
			this.$emit('update:adding', value);
			if (value) {
				this.$nextTick(this.focusInput);
			}
		},
	},
	methods: {
		focusInput() {
			// focus first input element
			const input = this.$el.querySelector('input');
			if (input) {
				input.focus();
			}
		},
		addTag() {
			if (this.addTagDisabled) {
				return;
			}
			ajaxPost('/ajax/space/create/tag', {
				parentId: this.parentId,
				tag: this.tag,
			}).then(response => {
				this.$emit('added', response);
				this.adding = false;
				this.tag = '';
			});
		},
	},
};
</script>
