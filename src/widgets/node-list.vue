<template>
<div class="node-list flex-column-md">

	<template v-for="n in nodes" :key="n.id">

		<node-header :node="n" link-to/>

	</template>

</div>
</template>

<script>
import NodeHeader from '@/widgets/node-header.vue';

export default {
	components: {
		NodeHeader,
	},
	props: {
		nodes: {
			type: Array,
			required: true,
		},
		parentId: {
			type: Number,
			required: false,
		},
	},
};
</script>
