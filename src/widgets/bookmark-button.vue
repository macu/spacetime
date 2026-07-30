<template>
<el-button @click.stop="toggleBookmark()" :type="buttonType" :size="size">
	<material-icon v-if="isBookmarked" icon="bookmark_border"/>
	<material-icon v-else icon="bookmark"/>
</el-button>
</template>

<script>
import {
	ajaxPost,
} from '@/utils/ajax.js';

import bus from '@/utils/bus.js';

export default {
	emits: [
		'bookmark-added',
		'bookmark-removed',
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
			isBookmarked: this.space.userBookmark || false,
		};
	},
	computed: {
		buttonType() {
			return this.isBookmarked ? 'success' : 'default';
		},
	},
	mounted() {
		bus.on('bookmark', this.setBookmarked);
	},
	beforeUnmount() {
		bus.off('bookmark', this.setBookmarked);
	},
	methods: {
		toggleBookmark() {
			let newState = !this.isBookmarked;
			let update = () => {
				ajaxPost('/ajax/bookmark', {
					spaceId: this.space.id,
					bookmark: newState,
				}).then(() => {
					if (newState) {
						this.$emit('bookmark-added');
					} else {
						this.$emit('bookmark-removed');
					}
					bus.emit('bookmark', {
						spaceId: this.space.id,
						newState,
					});
				});
			};
			if (newState) {
				update();
			} else {
				this.$confirm('Are you sure you want to remove this bookmark?',
					'Confirm remove bookmark', {
					type: 'warning',
					confirmButtonText: 'Yes',
					cancelButtonText: 'No',
				}).then(() => {
					update();
				}).catch(() => {
					// Cancelled
				});
			}
		},
		setBookmarked({spaceId, newState}) {
			if (this.space.id === spaceId) {
				this.isBookmarked = newState;
			}
		},
	},
};
</script>
