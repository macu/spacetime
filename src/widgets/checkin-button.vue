<template>
<el-button-group>
	<el-button @click="addCheckIn()" :type="buttonType">
		<material-icon icon="check"/>
	</el-button>
	<el-button @click="showStats()" :type="buttonType">
		<span v-text="checkinCount"/>
	</el-button>
</el-button-group>
</template>

<script>
export default {
	emits: [
		'check-in',
	],
	props: {
		space: {
			type: Object,
			required: true,
		},
	},
	data() {
		return {
			hasUserCheckin: false,
			checkinCount: this.space.checkinTotal || 0,
		};
	},
	computed: {
		buttonType() {
			return this.hasUserCheckin ? 'success' : 'primary';
		},
	},
	methods: {
		addCheckIn() {
			// TODO
		},
		showStats() {
			// TODO
		},
	},
};
</script>
