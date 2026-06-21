<template>

	<el-dropdown @command="mode => filter.mode = mode"
		trigger="click" placement="bottom-start">
		<el-button>
			<material-icon icon="filter_list"/>
			<span v-text="showingLabel"/>
			<material-icon icon="arrow_drop_down"/>
		</el-button>
		<template #dropdown>
			<el-dropdown-item v-if="filter.mode != FILTER_MODES.TOP" :command="FILTER_MODES.TOP">
				<material-icon icon="arrow_upward"/>
				<span>Top all time</span>
			</el-dropdown-item>
			<el-dropdown-item v-if="filter.mode != FILTER_MODES.RECENT" :command="FILTER_MODES.RECENT">
				<material-icon icon="history"/>
				<span>Most recent</span>
			</el-dropdown-item>
			<el-dropdown-item v-if="allowPinned && filter.mode != FILTER_MODES.PINNED" :command="FILTER_MODES.PINNED">
				<material-icon icon="push_pin"/>
				<span>Pinned</span>
			</el-dropdown-item>
		</template>
	</el-dropdown>

</template>

<script>
export const FILTER_MODES = {
	TOP: 'top-subspaces',
	RECENT: 'most-recent',
	PINNED: 'pinned',
};

export function getFilter() {
	return {
		mode: FILTER_MODES.TOP,
		date: null,
		window: null,
	};
}

export default {
	emits: [
		'update:modelValue',
	],
	props: {
		modelValue: {
			type: Object,
			default: () => getFilter(),
		},
		allowPinned: {
			// TODO pass in
			type: Boolean,
			default: true,
		},
	},
	data() {
		return {
			filter: this.modelValue || getFilter(),
		};
	},
	computed: {
		FILTER_MODES() {
			return FILTER_MODES;
		},
		showingLabel() {
			switch (this.filter.mode) {
				case FILTER_MODES.TOP:
					return 'Showing top all-time';
				case FILTER_MODES.RECENT:
					return 'Showing most recent';
				case FILTER_MODES.PINNED:
					return 'Showing pinned';
			}
			return 'Filter';
		},
	},
	watch: {
		filter: {
			deep: true,
			handler() {
				this.$emit('update:modelValue', this.filter);
			},
		},
	},
};
</script>
