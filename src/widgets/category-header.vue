<template>
<div class="category-header flex-column tree-node-border">
	<h3 class="flex-row-md">
		<material-icon :icon="icon"/>
		<router-link
			v-if="linkHeader"
			:to="{name: 'category-view', params: {id: node.id}}">
			<span v-text="node.title"/>
		</router-link>
		<span v-else v-text="node.title"/>
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
		linkHeader: {
			type: Boolean,
			required: false,
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
