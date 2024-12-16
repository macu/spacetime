<template>
<el-button-group :type="buttonType" :size="size">
	<el-button @click="addCheckIn()">
		<material-icon icon="check"/>
	</el-button>
	<el-button @click="showStats()">
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
		size: {
			type: String,
			default: null,
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
