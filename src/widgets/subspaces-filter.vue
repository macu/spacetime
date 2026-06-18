<template>

	<el-dropdown @command="mode => showing = mode"
		trigger="click" placement="bottom-start">
		<el-button>
			<material-icon icon="filter_list"/>
			<span v-text="showingLabel"/>
			<material-icon icon="arrow_drop_down"/>
		</el-button>
		<template #dropdown>
			<el-dropdown-item v-if="showing != MODES.TOP" :command="MODES.TOP">
				<material-icon icon="arrow_upward"/>
				<span>Top all time</span>
			</el-dropdown-item>
			<el-dropdown-item v-if="showing != MODES.RECENT" :command="MODES.RECENT">
				<material-icon icon="history"/>
				<span>Most recent</span>
			</el-dropdown-item>
			<el-dropdown-item v-if="showing != MODES.PINNED" :command="MODES.PINNED">
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

export default {
	props: {
		modelValue: {
			type: Object,
			default: null,
		},
	},
	data() {
		return {
			showing: this.modelValue?.mode || MODES.TOP,
		};
	},
	computed: {
		MODES() {
			return MODES;
		},
		showingLabel() {
			if (this.showing === MODES.TOP) {
				return 'Showing top all-time';
			} else if (this.showing === MODES.RECENT) {
				return 'Showing most recent';
			} else if (this.showing === MODES.PINNED) {
				return 'Showing pinned';
			}
			return 'Filter';
		},
	},
	watch: {
		showing: {
			immediate: true,
			handler() {
				this.$emit('update:filter', {
					mode: this.showing,
					date: null,
					window: null,
				});
			},
		},
	},
};
</script>
