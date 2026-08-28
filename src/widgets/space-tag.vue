<template>
<div class="space-tag flex-row nowrap" :class="{'ellipsis': ellipsis}">

	<template v-if="flat">
		<material-icon :icon="icon"/>
		<span
			@click="$emit('click-tag')"
			class="text"
			v-text="tagOutput"
			/>
	</template>

	<template v-else>

		<material-icon v-if="actionsExpanded" :icon="icon"/>
		<el-button v-else size="small" @click="actionsExpanded = true">
			<material-icon :icon="icon"/>
			<span v-text="space.voteSum"/>
		</el-button>

		<template v-if="actionsExpanded">

			<el-button
				v-if="showPinning"
				:type="pinned ? 'success' : 'primary'"
				size="small"
				@click="togglePinned()">
				<material-icon v-if="pinned" icon="keep"/>
				<material-icon v-else icon="keep_off"/>
				{{ pinned ? 'Unpin' : 'Pin' }}
			</el-button>

			<el-tooltip v-else-if="pinned" content="Pinned by author" placement="top">
				<material-icon icon="keep" class="pinned"/>
			</el-tooltip>

			<vote-actions
				v-if="showVoteActions"
				:space="space"
				@vote="$emit('vote')"
				size="small"
				/>

		</template>

		<el-tooltip v-else-if="pinned" content="Pinned by author" placement="top">
			<material-icon icon="keep" class="pinned"/>
		</el-tooltip>

		<span
			@click="$emit('click-tag')"
			class="text"
			v-text="tagOutput"
			/>

	</template>

</div>
</template>

<script>
import VoteActions from './vote-actions.vue';

import {
	SPACE_TYPE_ICONS,
} from '@/const.js';

export default {
	emits: [
		'vote',
		'click-tag',
		'toggle-pinned',
	],
	components: {
		VoteActions,
	},
	props: {
		space: {
			type: Object,
			required: true,
		},
		showVoteActions: {
			type: Boolean,
			default: true,
		},
		showPinning: {
			type: Boolean,
			default: false,
		},
		flat: {
			type: Boolean,
			default: false,
		},
		ellipsis: {
			type: Boolean,
			default: false,
		},
	},
	data() {
		return {
			actionsExpanded: false,
		};
	},
	computed: {
		icon() {
			return SPACE_TYPE_ICONS[this.space.spaceType];
		},
		tagOutput() {
			return this.space.text || '';
		},
		pinned() {
			return this.space.isPinned || false;
		},
	},
	methods: {
		togglePinned() {
			this.$emit('toggle-pinned');
		},
	},
};
</script>

<style lang="scss">
@import '@/styles/vars.scss';

.space-tag {
	padding: 0 0 5px;
	background-color: $tag-bg-color;
	border-bottom: $tag-border;
	color: $tag-color;
	overflow: hidden;
	cursor: pointer;

	>.text {
		font-size: 120%;
	}

	&.ellipsis {
		>.text {
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}
	}
}
</style>
