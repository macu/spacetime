<template>
<div ref="container" class="add-space-tag-widget">
	<div v-if="adding" class="add-tag-form flex-row-sm nowrap">
		<el-input
			v-model="tag"
			:maxlength="$const.tagMaxLength"
			show-word-limit
			size="small">
			<template #prepend>
				Add tag
			</template>
		</el-input>
		<el-button
			@click="addTag()"
			:disabled="addTagDisabled"
			size="small" type="primary">
			<material-icon icon="check"/>
		</el-button>
		<el-button
			@click="adding = false"
			size="small" type="warning">
			<material-icon icon="close"/>
		</el-button>
	</div>
	<el-button v-else @click="adding = true" type="primary" size="small" :disabled="disabled" plain>
		<material-icon icon="add"/>
		<span>Add tag</span>
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
			type: [Number, String],
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
			const input = this.$refs.container.querySelector('input');
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
			}, {
				409: 'Tag already added.', // conflict: already exists
			}).then(response => {
				this.$emit('added', response);
				this.adding = false;
				this.tag = '';
			});
		},
	},
};
</script>

<style lang="scss">
.add-space-tag-widget {
	.add-tag-form {
		border: thin solid lightgray;
		border-radius: 4px;
		padding: 5px;
	}
}
</style>
