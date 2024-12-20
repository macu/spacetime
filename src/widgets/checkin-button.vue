<template>
<el-button-group :type="buttonType" :size="size" class="checkin-button">
	<el-button @click="addCheckIn()">
		<material-icon icon="check"/>
	</el-button>
	<el-button v-if="totalSubspaces > 0" @click="showStats()">
		<span v-text="totalSubspaces"/>
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
			totalSubspaces: this.space.totalSubspaces || 0,
		};
	},
	computed: {
		buttonType() {
			return this.hasUserCheckin ? 'success' : 'primary';
		},
	},
	mounted() {
		bus.on('direct-check-in', this.incrementSubspaces);
	},
	beforeUnmount() {
		bus.off('direct-check-in', this.incrementSubspaces);
	},
	methods: {
		incrementSubspaces({spaceId}) {
			if (this.space.id === spaceId) {
				this.hasUserCheckin = true;
				this.totalSubspaces++;
			}
		},
		addCheckIn() {
			ajaxPost('/ajax/space/create/checkin', {
				parentId: this.space.id,
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
