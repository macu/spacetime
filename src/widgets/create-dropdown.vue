<template>
<div class="create-dropdown" :class="{'sticky': sticky}">
	<div>
		<el-dropdown v-if="hasParentId"
			@command="create" :disabled="disabled" placement="bottom-start">
			<el-button type="primary" :disabled="disabled">
				<material-icon icon="add"/>
				<span>Create</span>
				<material-icon icon="arrow_drop_down"/>
			</el-button>
			<template #dropdown>
				<el-dropdown-item command="create-branch">
					<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.BRANCH]"/>
					<span>Create branch</span>
				</el-dropdown-item>
				<template v-if="hasParentId">
					<el-dropdown-item command="create-space-link">
						<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.LINK]"/>
						<span>Create space link</span>
					</el-dropdown-item>
				</template>
				<el-dropdown-item command="create-text">
					<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.TEXT]"/>
					<span>Create text</span>
				</el-dropdown-item>
				<el-dropdown-item disabled command="create-stream-oc">
					<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.STREAM_OC]"/>
					<span>Create stream of consciousness</span>
				</el-dropdown-item>
				<el-dropdown-item disabled command="create-json-ar">
					<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.JSON_AR]"/>
					<span>Create json attribute reporter</span>
				</el-dropdown-item>
			</template>
		</el-dropdown>
		<el-button v-else
			@click="create('create-branch')"
			type="primary" :disabled="disabled">
			<material-icon icon="add"/>
			<span>Create space</span>
		</el-button>
	</div>
</div>
</template>

<script>
import {
	SPACE_TYPES,
	SPACE_TYPE_ICONS,
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
		SPACE_TYPE_ICONS() {
			return SPACE_TYPE_ICONS;
		},
		hasParentId() {
			return !!this.parentId;
		},
	},
	methods: {
		create(routeName) {
			let query = {};
			if (routeName.includes('?')) {
				let queryString;
				[routeName, queryString] = routeName.split('?');
				const params = new URLSearchParams(queryString);
				for (const [key, value] of params.entries()) {
					query[key] = value;
				}
			}
			this.$router.push({
				name: routeName,
				query: {
					parentId: this.parentId,
					...query,
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
