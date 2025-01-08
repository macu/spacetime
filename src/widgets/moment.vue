<template>
<el-tooltip effect="dark" :content="tooltipString" placement="bottom-start">
	<span class="moment" v-text="displayString"/>
</el-tooltip>
</template>

<script>
import moment from 'moment';

export default {
	props: {
		time: {
			type: String,
			required: true,
		},
	},
	computed: {
		displayString() {
			return moment(this.time).fromNow();
		},
		tooltipString() {
			return moment(this.time).format('YYYY-MM-DD HH:mm:ss');
		},
	},
};
</script>
