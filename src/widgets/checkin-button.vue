<template>
<el-button-group :type="buttonType" :size="size" class="checkin-button">
	<el-button @click="addCheckIn()" :disabled="disabled">
		<material-icon icon="check"/>
	</el-button>
	<el-button v-if="checkinCount > 0" @click="showStats()">
		<span v-text="checkinCount"/>
	</el-button>
</el-button-group>
</template>

<script>
import {
	showError,
} from '@/utils/notify.js';

import bus from '@/utils/bus.js';

import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	emits: [
		'check-in',
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
			hasUserCheckin: false,
			checkinCount: this.space.checkinCount || 0,
		};
	},
	computed: {
		buttonType() {
			return this.hasUserCheckin ? 'success' : 'primary';
		},
		disabled() {
			return this.$store.getters.createDisabled;
		},
	},
	mounted() {
		bus.on('direct-check-in', this.incrementCheckins);
	},
	beforeUnmount() {
		bus.off('direct-check-in', this.incrementCheckins);
	},
	methods: {
		incrementCheckins({spaceId}) {
			if (this.space.id === spaceId) {
				this.hasUserCheckin = true;
				this.checkinCount++;
			}
		},
		addCheckIn() {
			ajaxPost('/ajax/create/checkin', {
				parentId: this.space.id,
			}, {
				429() { // rate limit
					showError('Rate limit exceeded. Max 1 check-in per space per minute.');
				},
			}).then(() => {
				bus.emit('direct-check-in', {
					spaceId: this.space.id,
				});
				this.$emit('check-in');
			});
		},
		showStats() {
			// TODO
		},
	},
};
</script>

<style lang="scss">
.checkin-button {
	display: flex;
	flex-direction: row;
	flex-wrap: nowrap;
	>*, >*+* {
		margin: 0;
	}
}
</style>
