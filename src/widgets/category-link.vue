<template>
<div class="category-link flex-column tree-node-border">
	<h3 class="flex-row-md">
		<material-icon :icon="icon"/>
		<router-link :to="{name: 'category-view', params: {id: node.id}}">
			<span v-text="node.title"/>
		</router-link>
	</h3>
</div>
</template>

<script>
export default {
	props: {
		node: {
			type: Object,
			required: true,
		},
	},
	computed: {
		icon() {
			if (this.node.key) {
				switch (this.node.key) {
					case 'treetime':
						return 'park';
				}
			}
			return 'folder';
		},
	},
};
</script>
