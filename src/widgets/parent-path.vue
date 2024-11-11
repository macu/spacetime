<template>
<div class="node-parent-path">

	<div v-for="n in pathWithTitles" :key="n.id" class="node-parent flex-row-md nowrap">

		<node-icon :node="n"/>

		<router-link :to="{name: 'node-view', params: {id: n.id}}">
			<span v-text="n.title"/>
		</router-link>

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
		path: {
			type: Array,
			required: true,
		},
	},
	computed: {
		pathWithTitles() {
			return this.path.map(n => {
				switch (n.class) {
					case NODE_CLASS.LANG:
					case NODE_CLASS.TAG:
					case NODE_CLASS.CATEGORY:
					case NODE_CLASS.POST:
						return {
							...n,
							title: n.content.title,
						};
					case NODE_CLASS.COMMENT:
						return {
							...n,
							title: n.content.text,
						};
				}
			});
		},
	},
};
</script>

<style lang="scss">
@import '@/styles/vars.scss';

$border-width: 1px;

.node-parent-path {
	border: $border-width solid gray;
	background-color: gray;
	border-radius: $border-radius;
	box-shadow: $box-shadow;

	display: flex;
	flex-direction: column;
	row-gap: $border-width;

	>.node-parent {
		padding: 5px 20px;
		background-color: #eee;
		font-size: smaller;

		&:first-child {
			border-top-left-radius: $border-radius;
			border-top-right-radius: $border-radius;
		}
		&:last-child {
			border-bottom-left-radius: $border-radius;
			border-bottom-right-radius: $border-radius;
		}
	}
}
</style>
