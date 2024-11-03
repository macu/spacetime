<template>
<div class="category-view-page flex-column-lg page-width-md">

	<loading-message v-if="loadingNode"/>

	<template v-else-if="node">

		<div v-if="path.length" class="flex-column">
			<strong>Path</strong>
			<parent-path :path="path"/>
		</div>

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
			<template v-else-if="node.class === NODE_CLASS.TYPE">
				Type
			</template>
			<template v-else-if="node.class === NODE_CLASS.FIELD">
				Field
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
			<el-dropdown @command="gotoCreate" placement="bottom-start">
				<el-button type="primary">
					<span>Add subcontent</span>
					<material-icon icon="arrow_drop_down"/>
				</el-button>
				<template #dropdown>
					<el-dropdown-menu>
						<el-dropdown-item command="create-category">
							<material-icon icon="folder"/>
							<span>Category</span>
						</el-dropdown-item>
						<el-dropdown-item command="create-post">
							<material-icon icon="description"/>
							<span>Post</span>
						</el-dropdown-item>
					</el-dropdown-menu>
				</template>
			</el-dropdown>
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
import ParentPath from '@/widgets/parent-path.vue';
import NodeHeader from '@/widgets/node-header.vue';
import NodeList from '@/widgets/node-list.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

import {
	getPathKeyScope,
} from '@/utils/tree.js';

import {
	NODE_CLASS,
	SYSTEM_NODE_KEYS,
} from '@/const.js';

export default {
	components: {
		ParentPath,
		NodeHeader,
		NodeList,
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
		keyScope() {
			return getPathKeyScope([...(this.path || []), this.node]);
		},
		depth() {
			return this.path.length + 1;
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
		this.path = [];
		this.children = [];

		next();

		this.loadNode(to.params.id);
	},
	methods: {
		loadNode(id) {
			this.loadingNode = true;
			ajaxGet('/ajax/node/view', {
				id,
			}).then(response => {
				this.node = response.node;
				this.path = response.path;
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
		gotoCreate(routeName) {
			this.$router.push({
				name: routeName,
				query: {
					parentId: this.node.id,
				},
			});
		},
	},
};
</script>
