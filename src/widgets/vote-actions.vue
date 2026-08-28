<template>
<el-button v-if="disabled" :size="size" @click="showStats()" :type="buttonType">
	<material-icon v-if="voteSum >= 0" icon="thumb_up"/>
	<material-icon v-else icon="thumb_down"/>
	<span v-text="voteSum"/>
</el-button>
<el-button-group v-else :size="size" class="vote-actions">
	<el-button v-if="showUpvote" @click="toggleUpvote()"
		:type="currentVote === 1 ? 'success' : 'default'">
		<material-icon icon="thumb_up"/>
	</el-button>
	<el-button v-if="showDownvote" @click="toggleDownvote()"
		:type="currentVote === -1 ? 'danger' : 'default'">
		<material-icon icon="thumb_down"/>
	</el-button>
	<el-button @click="showStats()" :type="buttonType">
		<span v-text="voteSum"/>
	</el-button>
</el-button-group>
</template>

<script>
import bus from '@/utils/bus.js';

import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	emits: [
		'vote',
	],
	props: {
		space: {
			type: Object,
			required: true,
		},
		size: {
			type: String,
			default: null,
		},
	},
	data() {
		return {
			currentVote: this.space.currentVote,
			voteSum: this.space.voteSum,
		};
	},
	computed: {
		showUpvote() {
			return !this.disabled && (!this.currentVote || this.currentVote === 1);
		},
		showDownvote() {
			return !this.disabled && (!this.currentVote || this.currentVote === -1);
		},
		buttonType() {
			if (this.voteSum < 0) {
				return 'danger';
			} else if (this.voteSum > 0) {
				return 'success';
			}
			return 'default';
		},
		disabled() {
			// User not logged in
			return this.$store.getters.createDisabled;
		},
	},
	mounted() {
		bus.on('direct-set-vote', this.updateVote);
	},
	beforeUnmount() {
		bus.off('direct-set-vote', this.updateVote);
	},
	methods: {
		updateVote({ spaceId, currentVote, voteSum }) {
			if (spaceId === this.space.id) {
				console.log('Updating vote for space:', spaceId, 'Current vote:', currentVote, 'Vote sum:', voteSum);
				this.currentVote = currentVote;
				this.voteSum = voteSum;
			}
		},

		toggleUpvote() {
			if (this.currentVote === 1) {
				this.postVote(0);
			} else {
				this.postVote(1);
			}
		},
		toggleDownvote() {
			if (this.currentVote === -1) {
				this.postVote(0);
			} else {
				this.postVote(-1);
			}
		},

		postVote(voteValue) {
			if (this.disabled) {
				return;
			}
			ajaxPost('/ajax/space/vote', {
				spaceId: this.space.id,
				voteValue,
			}).then(response => {
				this.$emit('vote', voteValue);
				bus.emit('direct-set-vote', {
					spaceId: this.space.id,
					voteSum: response.voteSum,
					currentVote: voteValue,
				});
			});
		},

		showStats() {
			// TODO
		},
	},
};
</script>

<style lang="scss">
.vote-actions {
	display: flex;
	flex-direction: row;
	flex-wrap: nowrap;
	>*, >*+* {
		margin: 0;
	}
}
</style>
