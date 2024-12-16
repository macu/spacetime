<template>
<div>
	<el-input
		v-if="adding"
		v-model="title"
		:maxlength="$store.getters.titleMaxLength"
		show-word-limit>
		<template #prepend>
			Add title
		</template>
		<template #append>
			<el-button-group>
				<el-button
					@click="addTitle()"
					:disabled="addTitleDisabled"
					type="success">
					Add
				</el-button>
				<el-button
					@click="adding=false"
					type="warning">
					Cancel
				</el-button>
			</el-button-group>
		</template>
	</el-input>
	<el-button v-else @click="adding=true" type="primary">Add title</el-button>
</div>
</template>

<script>
import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	emits: [
		'added',
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
	},
	methods: {
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
