<template>
<div class="node-header flex-column">
	<h3 class="flex-row-md" :class="{'link-to': linkTo}" @click="gotoNode()">
		<node-icon :node="node"/>
		<router-link v-if="linkTo" :to="route">
			<span v-if="title" v-text="title"/>
			<em v-else-if="metaTitle" v-text="metaTitle"/>
			<em v-else>Untitled</em>
		</router-link>
		<span v-else-if="title" v-text="title"/>
		<em v-else-if="metaTitle" v-text="metaTitle"/>
		<em v-else>Untitled</em>
	</h3>
	<div v-if="description" v-text="description" class="pre-wrap"/>
	<div v-if="postBlocks" class="flex-column">
		<div v-for="b in postBlocks" v-text="b.text" class="pre-wrap"/>
		<div v-if="showExpand" class="center">
			<el-button @click="expanded = true">Show all</el-button>
		</div>
	</div>
	<div v-if="showActions" class="flex-row-md">
		<slot name="node-actions"/>
	</div>
</div>
</template>

<script>
import NodeIcon from '@/widgets/node-icon.vue';

import {
	NODE_CLASS,
} from '@/const.js';

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
		showAll: {
			type: Boolean,
			required: false,
		},
	},
	data() {
		return {
			expanded: false,
		};
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
		title() {
			switch (this.node.class) {
				case NODE_CLASS.LANG:
				case NODE_CLASS.TAG:
				case NODE_CLASS.CATEGORY:
				case NODE_CLASS.TYPE:
				case NODE_CLASS.FIELD:
				case NODE_CLASS.POST:
					return this.node.content.title;
			}
			return null;
		},
		metaTitle() {
			switch (this.node.class) {
				case NODE_CLASS.COMMENT:
					return 'Comment';
			}
			return null;
		},
		description() {
			return this.node.content.description || null;
		},
		showingAll() {
			return this.showAll || this.expanded;
		},
		postBlocks() {
			if (this.node.class === NODE_CLASS.POST) {
				if (this.showingAll) {
					return this.node.content.blocks;
				}
				return this.node.content.blocks.slice(0, 1);
			}
			return null;
		},
		showExpand() {
			return this.node.class === NODE_CLASS.POST &&
				!this.showingAll &&
				this.node.content.blocks.length > 1;
		},
		showActions() {
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
