<template>
<el-dropdown @command="create" :disabled="disabled" placement="bottom-start">
	<el-button type="primary" :disabled="disabled">
		<material-icon icon="add"/>
		<span>Create</span>
		<material-icon icon="arrow_drop_down"/>
	</el-button>
	<template #dropdown>
		<el-dropdown-item command="create-category">
			<material-icon icon="folder"/>
			<span>Create category</span>
		</el-dropdown-item>
		<el-dropdown-item command="create-post">
			<material-icon icon="description"/>
			<span>Create post</span>
		</el-dropdown-item>
		<el-dropdown-item command="create-comment">
			<material-icon icon="comment"/>
			<span>Create comment</span>
		</el-dropdown-item>
	</template>
</el-dropdown>
</template>

<script>
export default {
	props: {
		parentId: {
			type: Number,
			required: false,
		},
		disabled: {
			type: Boolean,
			default: false,
		},
	},
	methods: {
		create(routeName) {
			this.$router.push({
				name: routeName,
				query: {
					parentId: this.parentId,
				},
			});
		},
	},
};
</script>
