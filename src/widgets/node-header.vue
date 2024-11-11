<template>
<div class="node-header-wrapper">
	<parent-path v-if="hasPath" :path="path"/>
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
		<div v-else-if="commentText" v-text="commentText" class="pre-wrap"/>
		<div v-else-if="postBlocks" class="flex-column">
			<div v-for="b in postBlocks" v-text="b.text" class="pre-wrap"/>
		</div>
		<div v-if="showCreatorDetails" class="owner-info flex-row">
			<material-icon icon="person"/>
			<span>Created by {{node.creator.displayName}}</span>
		</div>
		<div v-if="showActions" class="flex-row-md">
			<slot name="node-actions"/>
		</div>
	</div>
</div>
</template>

<script>
import ParentPath from '@/widgets/parent-path.vue';
import NodeIcon from '@/widgets/node-icon.vue';

import {
	NODE_CLASS,
	OWNER_TYPE,
} from '@/const.js';

export default {
	components: {
		ParentPath,
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
	computed: {
		hasPath() {
			return !!this.node.path && this.node.path.length > 0;
		},
		path() {
			return this.node.path || null;
		},
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
				case NODE_CLASS.POST:
					return this.node.content.title;
			}
			return null;
		},
		metaTitle() {
			switch (this.node.class) {
				case NODE_CLASS.COMMENT:
					return 'Comment by ' + this.node.creator.displayName;
			}
			return null;
		},
		description() {
			return this.node.content.description || null;
		},
		commentText() {
			if (this.node.class === NODE_CLASS.COMMENT) {
				return this.node.content.text;
			}
			return null;
		},
		postBlocks() {
			if (this.node.class === NODE_CLASS.POST) {
				return this.node.content.blocks;
			}
			return null;
		},
		showCreatorDetails() {
			return this.node.ownerType === OWNER_TYPE.USER;
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

$border-width: 2px;

.node-header-wrapper {
	>.node-parent-path {
		border-bottom: none;
		border-bottom-left-radius: 0;
		border-bottom-right-radius: 0;
		>.node-parent:last-child {
			border-bottom-left-radius: 0;
			border-bottom-right-radius: 0;
		}
	}
	>.node-header {
		padding: 20px;
		border: $border-width solid gray;
		border-radius: $border-radius;
		box-shadow: $box-shadow;

		>.link-to {
			cursor: pointer;
		}

		>.owner-info {
			font-size: smaller;
		}
	}
	>.node-parent-path+.node-header {
		border-top-left-radius: 0;
		border-top-right-radius: 0;
	}
}
</style>
