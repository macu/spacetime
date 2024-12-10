<template>
<el-button-group>
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
	props: {
		space: {
			type: Object,
			required: true,
		},
	},
	data() {
		return {
			hasUserCheckin: false,
			checkinCount: this.space.checkinCountAll || 0,
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
