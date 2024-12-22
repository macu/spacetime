<template>
<div class="parent-path">

	<loading-message v-if="loading"/>

	<space-output
		v-if="parentSpace"
		:space="parentSpace"
		show-path
		/>

</div>
</template>

<script>
import SpaceOutput from '@/widgets/space-output.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		SpaceOutput,
	},
	props: {
		parentId: {
			type: [String, Number],
			required: true,
		},
	},
	data() {
		return {
			loading: true,
			parentSpace: null,
		};
	},
	mounted() {
		ajaxGet('/ajax/space', {
			spaceId: this.parentId,
			includeParentPath: true,
		}).then(response => {
			this.parentSpace = response;
		}).finally(() => {
			this.loading = false;
		});
	},
};
</script>
