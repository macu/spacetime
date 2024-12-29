<template>
<div class="add-space-title-widget">
	<div v-if="adding" class="flex-row nowrap">
		<el-input
			v-model="title"
			:maxlength="$store.getters.titleMaxLength"
			show-word-limit
			size="large">
			<template #prepend>
				Add title
			</template>
		</el-input>
		<el-button
			@click="addTitle()"
			:disabled="addTitleDisabled"
			size="large" type="primary">
			<material-icon icon="check"/>
		</el-button>
		<el-button
			@click="adding = false"
			size="large" type="warning">
			<material-icon icon="close"/>
		</el-button>
	</div>
	<el-button v-else @click="adding = true" type="primary" :disabled="disabled" plain>
		<material-icon icon="add"/>
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
			// Title added under parent space
			type: Number,
			required: true,
		},
	},
	data() {
		return {
			adding: false,
			title: '',
		};
	},
	computed: {
		addTitleDisabled() {
			return !this.title.trim();
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
		addTitle() {
			if (this.addTitleDisabled) {
				return;
			}
			ajaxPost('/ajax/space/create/title', {
				parentId: this.parentId,
				title: this.title,
			}).then(response => {
				this.$emit('added', response);
				this.adding = false;
				this.title = '';
			});
		},
	},
};
</script>
