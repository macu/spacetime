<template>
<div class="parent-path">

	<div v-for="n in path" :key="n.id" class="flex-row-md nowrap">

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
		NODE_CLASS() {
			return NODE_CLASS;
		},
	},
};
</script>

<style lang="scss">
$border-radius: 10px;
$border-width: 2px;

.parent-path {
	border: $border-width solid gray;
	background-color: gray;
	border-radius: $border-radius;
	box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);

	display: flex;
	flex-direction: column;
	row-gap: $border-width;

	>div {
		padding: 10px 20px;
		background-color: #eee;

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
