<template>
<div class="category-view-page flex-column-lg page-width-md">

	<loading-message v-if="loadingNode"/>

	<template v-else-if="node">

		<h2>
			<template v-if="node.class === NODE_CLASS.CATEGORY">
				Category
			</template>
			<template v-else-if="node.class === NODE_CLASS.LANG">
				Language
			</template>
			<template v-else-if="node.class === NODE_CLASS.TAG">
				Tag
			</template>
			<template v-else-if="node.class === NODE_CLASS.POST">
				Post
			</template>
			<template v-else-if="node.class === NODE_CLASS.COMMENT">
				Comment
			</template>
			<template v-else>
				Node
			</template>
		</h2>

		<node-header :node="node" show-all/>

		<horizontal-controls>
			<create-dropdown :parent-id="node.id" :disabled="disableCreate"/>
		</horizontal-controls>

		<loading-message v-if="loadingChildren"/>

		<node-list
			v-else-if="children.length"
			:nodes="children"
			:parent-id="node.id"
			/>

		<el-alert v-else
			title="No subcontent currently exists."
			type="info"
			:closable="false"
			/>

	</template>

	<el-alert v-else
		title="Node could not be loaded."
		type="error" show-icon :closable="false"
		/>

</div>
</template>

<script>
import NodeHeader from '@/widgets/node-header.vue';
import NodeList from '@/widgets/node-list.vue';
import CreateDropdown from '@/widgets/create-dropdown.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

import {
	NODE_CLASS,
	SYSTEM_NODE_KEYS,
} from '@/const.js';

export default {
	components: {
		NodeHeader,
		NodeList,
		CreateDropdown,
	},
	data() {
		return {
			loadingNode: true,
			node: null,
			path: [],

			loadingChildren: true,
			children: [],
		};
	},
	computed: {
		id() {
			return this.$route.params.id;
		},
		NODE_CLASS() {
			return NODE_CLASS;
		},
		authenticated() {
			return this.$store.getters.userIsAuthenticated;
		},
		disableCreate() {
			return !this.authenticated;
		},
	},
	beforeRouteEnter(to, from, next) {
		next(vm => {
			vm.loadNode(to.params.id);
		});
	},
	beforeRouteUpdate(to, from, next) {
		this.loadingNode = true;
		this.node = null;
		this.children = [];

		next();

		this.loadNode(to.params.id);
	},
	methods: {
		loadNode(id) {
			this.loadingNode = true;
			ajaxGet('/ajax/node/view', {
				id,
			}).then(node => {
				this.node = node;
				this.loadChildren();
			}).finally(() => {
				this.loadingNode = false;
			});
		},
		loadChildren() {
			if (!this.node) {
				return;
			}
			this.loadingChildren = true;
			ajaxGet('/ajax/node/children', {
				id: this.node.id,
			}).then(response => {
				this.children = response.nodes;
			}).finally(() => {
				this.loadingChildren = false;
			});
		},
		gotoCreateCategory() {
			this.$router.push({
				name: 'create-category',
				query: {parentId: this.node.id},
			});
		},
		gotoCreatePost() {
			this.$router.push({
				name: 'create-post',
				query: {parentId: this.node.id},
			});
		},
	},
};
</script>
