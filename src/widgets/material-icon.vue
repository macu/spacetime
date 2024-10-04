<template>
<span class="material-symbols-outlined" :class="classes">{{icon}}</span>
</template>

<script>
export default {
	props: {
		icon: {
			type: String,
			required: true,
		},
		large: {
			type: Boolean,
			default: false,
		},
	},
	computed: {
		classes() {
			if (this.large) {
				return ['large'];
			}
			return [];
		},
	},
};
</script>

<style lang="scss">
.material-icon {
	font-family: 'Material Icons';
	font-size: 24px;
	&.large {
		font-size: 48px;
	}
}
</style>
