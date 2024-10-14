<template>
<div class="node-header flex-column">
	<h3 class="flex-row-md">
		<node-icon :node="node"/>
		<router-link v-if="linkTo" :to="route">
			<span v-text="node.title"/>
		</router-link>
		<span v-else v-text="node.title"/>
	</h3>
	<div v-if="showContent" class="flex-row-md">
		<slot name="node-actions"/>
	</div>
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
		showContent() {
			return !!this.$slots['node-actions'];
		},
	},
};
</script>

<style lang="scss">
@import '@/styles/vars.scss';

.node-header {
	padding: 20px;
	border: 2px solid gray;
	border-radius: $border-radius;
	box-shadow: $box-shadow;
}
</style>
