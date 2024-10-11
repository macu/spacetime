<template>
<div class="tag-category-body flex-column-lg">

	<el-alert
		v-if="depthExceeded"
		title="Maximum depth reached"
		type="warning"
		:closable="false"
		/>

	<div v-else class="flex-row-md">
		<el-button @click="addSubcategory()" type="primary">Create subcategory</el-button>
		<el-button @click="addTag()" type="primary">Create tag</el-button>
	</div>

	<loading-message v-if="loadingChildren"/>

	<node-list
		v-else-if="children.length"
		:nodes="children"
		:parent-id="node.id"
		/>

</div>
</template>

<script>
import NodeList from '@/widgets/node-list.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		NodeList,
	},
	props: {
		node: {
			type: Object,
			required: true
		},
		depth: {
			type: Number,
			required: true
		},
	},
	data() {
		return {
			loadingChildren: true,
			children: [],
		};
	},
	computed: {
		depthExceeded() {
			return this.depth >= this.$store.getters.treeMaxDepth;
		},
	},
	mounted() {
		this.loadChildren();
	},
	methods: {
		loadChildren() {
			this.loadingChildren = true;
			ajaxGet('/ajax/node/children', {
				id: this.node.id,
			}).then(response => {
				this.children = response.nodes;
			}).finally(() => {
				this.loadingChildren = false;
			});
		},
	},
};
</script>
