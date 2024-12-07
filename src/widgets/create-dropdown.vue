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
				<el-dropdown-item :command="SPACE_TYPES.SPACE">
					<material-icon icon="folder"/>
					<span>Create empty space</span>
				</el-dropdown-item>
				<el-dropdown-item v-if="hasParentId" :command="SPACE_TYPES.CHECK_IN">
					<material-icon icon="folder"/>
					<span>Create check-in</span>
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
	computed() {
		SPACE_TYPES() {
			return SPACE_TYPES;
		},
		hasParentId() {
			return !!this.parentId;
		},
	},
	methods: {
		create(spaceType) {
			this.$router.push({
				name: 'create-space',
				query: {
					parentId: this.parentId,
					type: spaceType,
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
