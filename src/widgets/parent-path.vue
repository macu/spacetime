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
};
</script>

<style lang="scss">
@import '@/styles/vars.scss';

$border-width: 1px;

.parent-path {
	border: $border-width solid gray;
	background-color: gray;
	border-radius: $border-radius;
	box-shadow: $box-shadow;

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
