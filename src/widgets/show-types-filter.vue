<template>
<div class="show-types-filter flex-row-md">
	<el-checkbox v-model="showAll">
		<template v-if="showAll">
			Show all subspace types
		</template>
		<template v-else>
			Show all
		</template>
	</el-checkbox>
	<el-checkbox-group v-model="showTypes">
		<el-checkbox :label="SPACE_TYPES.SPACE">
			<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.SPACE]"/>
			<span>Space</span>
		</el-checkbox>
		<el-checkbox :label="SPACE_TYPES.CHECK_IN">
			<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.CHECK_IN]"/>
			<span>Check in</span>
		</el-checkbox>
		<el-checkbox :label="SPACE_TYPES.TITLE">
			<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.TITLE]"/>
			<span>Title</span>
		</el-checkbox>
		<el-checkbox :label="SPACE_TYPES.TAG">
			<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.TAG]"/>
			<span>Tag</span>
		</el-checkbox>
		<el-checkbox :label="SPACE_TYPES.TEXT">
			<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.TEXT]"/>
			<span>Text</span>
		</el-checkbox>
		<el-checkbox :label="SPACE_TYPES.NAKED_TEXT">
			<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.NAKED_TEXT]"/>
			<span>Naked text</span>
		</el-checkbox>
		<el-checkbox :label="SPACE_TYPES.STREAM_OC">
			<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.STREAM_OC]"/>
			<span>Stream of consciousness</span>
		</el-checkbox>
		<el-checkbox :label="SPACE_TYPES.JSON_AR">
			<material-icon :icon="SPACE_TYPE_ICONS[SPACE_TYPES.JSON_AR]"/>
			<span>JSON attribute</span>
		</el-checkbox>
	</el-checkbox-group>
</div>
</template>

<script>
import {
	getStorage,
	setStorage,
} from '@/utils/storage.js';

import {
	SPACE_TYPES,
	SPACE_TYPE_ICONS,
} from '@/const.js';

export default {
	emits: ['update:showTypes'],
	data() {
		return {
			showAll: false,
			showTypes: getStorage('filterShowTypes', Object.values(SPACE_TYPES)),
		};
	},
	computed: {
		SPACE_TYPES() {
			return SPACE_TYPES;
		},
		selectedTypes() {
			return this.showAll ? Object.values(SPACE_TYPES) : this.showTypes;
		},
	},
	watch: {
		selectedTypes: {
			deep: true,
			immediate: true,
			handler(selected) {
				setStorage('filterShowTypes', selected);
				this.$emit('update:showTypes', selected);
			},
		},
	},
};
</script>
