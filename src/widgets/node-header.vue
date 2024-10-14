<template>
<div class="node-header flex-column">
	<h3 class="flex-row-md" :class="{'link-to': linkTo}" @click="gotoNode()">
		<node-icon :node="node"/>
		<router-link v-if="linkTo" :to="route">
			<span v-if="node.title" v-text="node.title"/>
			<em v-else>Untitled</em>
		</router-link>
		<span v-else-if="node.title" v-text="node.title"/>
		<em v-else>Untitled</em>
	</h3>
	<div v-if="node.body" v-text="node.body" class="flex-column"/>
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
	methods: {
		gotoNode() {
			if (this.linkTo) {
				this.$router.push(this.route);
			}
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

	>.link-to {
		cursor: pointer;
	}
}
</style>
