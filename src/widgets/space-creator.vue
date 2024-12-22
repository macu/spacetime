<template>
<span class="space-creator">
	Created by <router-link :to="route" v-text="creatorName"/> <moment :time="space.createdAt"/>
</span>
</template>

<script>
export default {
	props: {
		space: {
			type: Object,
			required: true,
		},
	},
	computed: {
		creatorName() {
			return this.space.authorDisplayName || 'Unknown';
		},
		route() {
			return {
				name: 'user',
				params: {
					id: this.space.authorHandle || this.space.createdBy,
				},
			};
		},
	},
};
</script>
