<template>
<span class="moment" v-text="displayString"/>
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
	},
};
</script>
