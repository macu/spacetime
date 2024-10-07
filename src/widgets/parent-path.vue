<template>
<div class="parent-path">
	<div v-for="n in path" :key="n.id" class="flex-row-md">

		<material-icon :icon="nodeIcon(n)"/>

		<router-link v-if="n.class === NODE_CLASS.CATEGORY"
			:to="{name: 'category-view', params: {id: n.id}}">
			<span v-text="n.title"/>
		</router-link>

	</div>
</div>
</template>

<script>
import {
	NODE_CLASS,
} from '@/const.js';

export default {
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
	methods: {
		nodeIcon(n) {
			if (n.key) {
				switch (n.key) {
					case 'treetime':
						return 'park';
				}
			}
			return 'subdirectory_arrow_right';
		},
	},
};
</script>

<style lang="scss">
$border-radius: 10px;

.parent-path {
	display: flex;
	flex-direction: column;

	border: thin solid gray;
	border-radius: $border-radius;
	box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);

	>div {
		padding: 20px;
		font-size: 120%;

		&:first-child {
			border-top-left-radius: $border-radius;
			border-top-right-radius: $border-radius;
		}
		&:last-child {
			border-bottom-left-radius: $border-radius;
			border-bottom-right-radius: $border-radius;
		}
		&:not(:first-child) {
			border-top: thin solid gray;
		}
	}
}
</style>
