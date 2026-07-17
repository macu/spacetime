<template>
<span class="space-creator">
	Created by <span @click.stop><router-link :to="route" v-text="creatorName"/></span> <moment :time="space.createdAt"/>
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
