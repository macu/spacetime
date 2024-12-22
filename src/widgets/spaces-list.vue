<template>
<div class="spaces-list flex-column-md">

	<space-output v-for="s in spaces" :space="s"/>

	<div class="center">
		<loading-message v-if="loading"/>
		<el-button v-else @click="$emit('load-more')" type="primary">Load More</el-button>
	</div>

</div>
</template>

<script>
import SpaceOutput from './space-output.vue';

export default {
	emits: [
		'load-more',
	],
	components: {
		SpaceOutput,
	},
	props: {
		spaces: {
			type: Array,
			required: true,
		},
		loading: {
			type: Boolean,
			default: false,
		},
	},
};
</script>
