<template>
<div class="post-body flex-column-lg">

	<el-alert
		v-if="depthExceeded"
		title="Maximum depth reached"
		type="warning"
		:closable="false"
		/>

	<div v-else class="flex-row-md">
		<el-button @click="addComment()" type="primary">Add comment</el-button>
	</div>

	<loading-message v-if="loadingComments" message="Loading comments..."/>

	<node-list
		v-else-if="comments.length"
		:nodes="comments"
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
			loadingComments: true,
			comments: [],
		};
	},
	computed: {
		depthExceeded() {
			return this.depth >= this.$store.getters.treeMaxDepth;
		},
	},
	mounted() {
		this.loadComments();
	},
	methods: {
		loadComments() {
			this.loadingComments = true;
			ajaxGet('/ajax/node/comments', {
				id: this.node.id,
			}).then(response => {
				this.comments = response.nodes;
			}).finally(() => {
				this.loadingComments = false;
			});
		},
	},
};
</script>
