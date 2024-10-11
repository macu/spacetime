<template>
<div class="node-header flex-column tree-node-border">
	<h3 class="flex-row-md">
		<node-icon :node="node"/>
		<router-link v-if="linkTo" :to="route">
			<span v-text="node.title"/>
		</router-link>
		<span v-else v-text="node.title"/>
	</h3>
</div>
</template>

<script>
import NodeIcon from '@/widgets/node-icon.vue';

export default {
	components: {
		NodeIcon,
	},
	props: {
		node: {
			type: Object,
			required: true,
		},
		linkTo: {
			type: Boolean,
			required: false,
		},
	},
	computed: {
		route() {
			return {
				name: 'node-view',
				params: {
					id: this.node.id,
				},
			};
		},
	},
};
</script>
