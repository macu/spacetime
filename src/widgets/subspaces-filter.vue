<template>

	<el-dropdown @command="mode => filter.mode = mode"
		trigger="click" placement="bottom-start">
		<el-button>
			<material-icon icon="filter_list"/>
			<span v-text="showingLabel"/>
			<material-icon icon="arrow_drop_down"/>
		</el-button>
		<template #dropdown>
			<el-dropdown-item v-if="filter.mode != MODES.TOP" :command="MODES.TOP">
				<material-icon icon="arrow_upward"/>
				<span>Top all time</span>
			</el-dropdown-item>
			<el-dropdown-item v-if="filter.mode != MODES.RECENT" :command="MODES.RECENT">
				<material-icon icon="history"/>
				<span>Most recent</span>
			</el-dropdown-item>
			<el-dropdown-item v-if="filter.mode != MODES.PINNED" :command="MODES.PINNED">
				<material-icon icon="push_pin"/>
				<span>Pinned</span>
			</el-dropdown-item>
		</template>
	</el-dropdown>

</template>

<script>
const MODES = {
	TOP: 'top-subspaces',
	RECENT: 'most-recent',
	PINNED: 'pinned',
};

export function getFilter() {
	return {
		mode: MODES.TOP,
		date: null,
		window: null,
	};
}

export default {
	emits: [
		'update:filter',
	],
	props: {
		modelValue: {
			type: Object,
			default: () => getFilter(),
		},
	},
	data() {
		return {
			filter: this.modelValue || getFilter(),
		};
	},
	computed: {
		MODES() {
			return MODES;
		},
		showingLabel() {
			switch (this.filter.mode) {
				case MODES.TOP:
					return 'Showing top all-time';
				case MODES.RECENT:
					return 'Showing most recent';
				case MODES.PINNED:
					return 'Showing pinned';
			}
			return 'Filter';
		},
	},
	watch: {
		filter: {
			deep: true,
			handler() {
				this.$emit('update:filter', this.filter);
			},
		},
	},
};
</script>
