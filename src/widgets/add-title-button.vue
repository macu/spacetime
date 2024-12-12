<template>
<el-input
	v-if="adding"
	v-model="title"
	:maxlength="$store.getters.titleMaxLength"
	show-word-limit>
	<template #prepend>
		Add title
	</template>
	<template #append>
		<el-button @click="addTitle()" :disabled="addTitleDisabled" type="success">Add</el-button>
		<el-button @click="adding=false" type="warning">Cancel</el-button>
	</template>
</el-input>
<el-button v-else @click="adding=true" type="primary">Add title</el-button>
</template>

<script>
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
			// TODO
		},
	},
};
</script>
