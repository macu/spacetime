<template>
<p class="loading-message flex-row nowrap" :class="align">
	<material-icon icon="loop" class="spin"/>
	<span v-text="displayMessage"/>
</p>
</template>

<script>
export default {
	props: {
		align: {
			type: String,
			default: 'left',
		},
		message: {
			type: String,
			required: false,
		},
	},
	computed: {
		displayMessage() {
			if (this.message) {
				return this.message;
			}
			return 'Loading...';
		},
	},
};
</script>

<style lang="scss">
.loading-message {
	font-size: larger;

	&.left {
		justify-content: flex-start;
	}
	&.center {
		justify-content: center;
	}
	&.right {
		justify-content: flex-end;
	}
}
</style>
