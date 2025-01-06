<template>

	<el-dropdown @command="show" placement="bottom-start" :hide-on-click="false">
		<el-button>
			<material-icon icon="filter_list"/>
			<span>Filter</span>
			<material-icon icon="arrow_drop_down"/>
		</el-button>
		<template #dropdown>
			<el-dropdown-item v-if="!showingAll" command="all">
				<material-icon icon="done_all"/>
				<span>Show all types</span>
			</el-dropdown-item>
			<el-dropdown-item
				v-for="type in SPACE_TYPES"
				:key="type"
				:command="type">
				<div class="show-types-filter-type flex-row"
					:class="{'showing': showTypes.includes(type)}">
					<material-icon :icon="SPACE_TYPE_ICONS[type]"/>
					<span v-text="titlesByType[type]"/>
					<material-icon v-if="showTypes.includes(type)" icon="check"/>
				</div>
			</el-dropdown-item>
		</template>
	</el-dropdown>

	<el-tag v-if="showingAll" size="large">
		<material-icon icon="done_all"/>
		<span>Showing all types</span>
	</el-tag>

	<template v-else>
		<el-tag v-for="t in showTypes" :key="t" closable @close="toggleShow(t)" size="large">
			<span class="flex-row">
				<material-icon :icon="SPACE_TYPE_ICONS[t]"/>
				<span v-text="titlesByType[t]"/>
			</span>
		</el-tag>
	</template>

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
			showTypes: getStorage('filterShowTypes', Object.values(SPACE_TYPES)),
		};
	},
	computed: {
		SPACE_TYPES() {
			return SPACE_TYPES;
		},
		titlesByType() {
			return {
				[SPACE_TYPES.SPACE]: 'Space',
				[SPACE_TYPES.SPACE_LINK]: 'Space link',
				[SPACE_TYPES.CHECK_IN]: 'Check in',
				[SPACE_TYPES.TITLE]: 'Title',
				[SPACE_TYPES.TAG]: 'Tag',
				[SPACE_TYPES.TEXT]: 'Text',
				[SPACE_TYPES.NAKED_TEXT]: 'Naked text',
				[SPACE_TYPES.STREAM_OC]: 'Stream of consciousness',
				[SPACE_TYPES.JSON_AR]: 'JSON attribute',
			};
		},
		SPACE_TYPE_ICONS() {
			return SPACE_TYPE_ICONS;
		},
		showingAll() {
			return Object.values(SPACE_TYPES).every(type => this.showTypes.includes(type));
		},
		typesNotShowing() {
			return Object.values(SPACE_TYPES).filter(type => !this.showTypes.includes(type));
		},
	},
	watch: {
		showTypes: {
			deep: true,
			immediate: true,
			handler(selected) {
				setStorage('filterShowTypes', selected);
				this.$emit('update:showTypes', selected);
			},
		},
	},
	methods: {
		showAll() {
			this.showTypes = Object.values(SPACE_TYPES);
		},
		toggleShow(type) {
			if (this.showTypes.includes(type)) {
				this.showTypes = this.showTypes.filter(t => t !== type);
			} else {
				this.showTypes.push(type);
			}
		},
		show(type) {
			if (type === 'all') {
				this.showAll();
			} else {
				this.toggleShow(type);
			}
		},
	},
};
</script>

<style lang="scss">
.show-types-filter-type {
	&.showing {
		font-weight: bold;
	}
}
</style>
