<template>
<div class="node-list flex-column-md">

	<template v-for="n in nodes" :key="n.id">

		<node-header :node="n" link-to>
			<template v-if="showNodeActions" #node-actions>
				<slot name="node-actions" :node="n"/>
			</template>
		</node-header>

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
		autoloadMore: {
			type: Boolean,
			required: false,
			default: false,
		},
	},
	computed: {
		showNodeActions() {
			return !!this.$slots['node-actions'];
		},
	},
};
</script>
