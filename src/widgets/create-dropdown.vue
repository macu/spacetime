<template>
<div class="create-dropdown" :class="{'sticky': sticky}">
	<div>
		<el-dropdown @command="create" :disabled="disabled" placement="bottom-start">
			<el-button type="primary" :disabled="disabled">
				<material-icon icon="add"/>
				<span>Create</span>
				<material-icon icon="arrow_drop_down"/>
			</el-button>
			<template #dropdown>
				<el-dropdown-item command="create-empty-space">
					<material-icon icon="folder"/>
					<span>Create empty space</span>
				</el-dropdown-item>
				<el-dropdown-item disabled command="create-stream-oc">
					<material-icon icon="folder"/>
					<span>Create stream of consciousness</span>
				</el-dropdown-item>
				<template v-if="hasParentId">
					<el-dropdown-item command="create-check-in">
						<material-icon icon="folder"/>
						<span>Create check-in</span>
					</el-dropdown-item>
					<el-dropdown-item command="create-title">
						<material-icon icon="folder"/>
						<span>Create title</span>
					</el-dropdown-item>
					<el-dropdown-item command="create-tag">
						<material-icon icon="folder"/>
						<span>Create tag</span>
					</el-dropdown-item>
					<el-dropdown-item command="create-text">
						<material-icon icon="folder"/>
						<span>Create text</span>
					</el-dropdown-item>
					<el-dropdown-item command="create-naked-text">
						<material-icon icon="folder"/>
						<span>Create naked text</span>
					</el-dropdown-item>
					<el-dropdown-item disabled command="create-json-ar">
						<material-icon icon="folder"/>
						<span>Create json attribute reporter</span>
					</el-dropdown-item>
				</template>
			</template>
		</el-dropdown>
	</div>
</div>
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
		sticky: {
			type: Boolean,
			default: false,
		},
	},
	computed: {
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

<style lang="scss">
.create-dropdown {
	&.sticky {
		position: sticky;
		top: 0;
		z-index: 100;

		display: flex;
		flex-direction: row;

		>div {
			align-self: flex-start;
			background-color: white;
			padding: 10px;
		}
	}
}
</style>
