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
				Language tag
			</template>
			<template v-else-if="node.class === NODE_CLASS.TAG">
				Tag
			</template>
			<template v-else-if="node.class === NODE_CLASS.TYPE">
				Type tag
			</template>
			<template v-else-if="node.class === NODE_CLASS.FIELD">
				Field tag
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

		<node-header :node="node"/>

		<component
			v-if="bodyClass"
			:is="bodyClass"
			:node="node"
			:depth="depth"
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

// Body classes
import Category from './category.vue';
import Lang from './lang.vue';
import LangCategory from './lang-category.vue';
import Tag from './tag.vue';
import TagCategory from './tag-category.vue';
import Type from './type.vue';
import TypeCategory from './type-category.vue';
import Field from './field.vue';
import Post from './post.vue';
import Comment from './comment.vue';

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
	},
	data() {
		return {
			loadingNode: true,
			node: null,
			path: [],
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
		bodyClass() {
			switch (this.node.class) {
				case NODE_CLASS.CATEGORY:
					if (this.keyScope === SYSTEM_NODE_KEYS.LANGS) {
						return LangCategory;
					} else if (this.keyScope === SYSTEM_NODE_KEYS.TAGS) {
						return TagCategory;
					} else if (this.keyScope === SYSTEM_NODE_KEYS.TYPES) {
						if (this.node.class === NODE_CLASS.TYPE) {
							return Type;
						} else {
							return TypeCategory;
						}
					}
					return Category;
				case NODE_CLASS.LANG:
					return Lang;
				case NODE_CLASS.TAG:
					return Tag;
				case NODE_CLASS.TYPE:
					return Type;
				case NODE_CLASS.FIELD:
					return Field;
				case NODE_CLASS.POST:
					return Post;
				case NODE_CLASS.COMMENT:
					return Comment;
			}
			return null;
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
			}).finally(() => {
				this.loadingNode = false;
			});
		},
	},
};
</script>
