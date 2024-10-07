<template>
<div class="nodes-list flex-column-md">

	<template v-for="n in nodes" :key="n.id">

		<category-header
			v-if="n.class === NODE_CLASS.CATEGORY" :node="n"
			link-header
			/>

	</template>

</div>
</template>

<script>
import CategoryHeader from '@/widgets/category-header.vue';

import {
	NODE_CLASS,
} from '@/const.js';

export default {
	components: {
		CategoryHeader,
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
	computed: {
		NODE_CLASS() {
			return NODE_CLASS;
		},
	},
};
</script>
