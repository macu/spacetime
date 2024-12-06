<template>
<el-dropdown @command="create" :disabled="disabled" placement="bottom-start">
	<el-button type="primary" :disabled="disabled">
		<material-icon icon="add"/>
		<span>Create</span>
		<material-icon icon="arrow_drop_down"/>
	</el-button>
	<template #dropdown>
		<el-dropdown-item :command="SPACE_TYPES.SPACE">
			<material-icon icon="folder"/>
			<span>Create empty space</span>
		</el-dropdown-item>
		<el-dropdown-item v-if="hasParentId" :command="SPACE_TYPES.CHECK_IN">
			<material-icon icon="folder"/>
			<span>Create check-in</span>
		</el-dropdown-item>
		<el-dropdown-item v-if="hasParentId" :command="SPACE_TYPES.TITLE">
			<material-icon icon="description"/>
			<span>Create title for this space</span>
		</el-dropdown-item>
		<el-dropdown-item v-if="hasParentId" :command="SPACE_TYPES.TAG">
			<material-icon icon="comment"/>
			<span>Create tag for this space</span>
		</el-dropdown-item>
		<el-dropdown-item :command="SPACE_TYPES.TEXT">
			<material-icon icon="folder"/>
			<span>Create text</span>
		</el-dropdown-item>
		<el-dropdown-item :command="SPACE_TYPES.NAKED_TEXT">
			<material-icon icon="folder"/>
			<span>Create naked text</span>
		</el-dropdown-item>
		<el-dropdown-item :command="SPACE_TYPES.STREAM_OC">
			<material-icon icon="folder"/>
			<span>Create stream of consciousness</span>
		</el-dropdown-item>
		<el-dropdown-item :command="SPACE_TYPES.JSON_AR">
			<material-icon icon="folder"/>
			<span>Create json attribute reporter</span>
		</el-dropdown-item>
	</template>
</el-dropdown>
</template>

<script>
import {
	SPACE_TYPES,
} from '@/const.js';

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
	computed() {
		SPACE_TYPES() {
			return SPACE_TYPES;
		},
		hasParentId() {
			return !!this.parentId;
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
