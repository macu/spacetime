<template>

	<loading-message v-if="loading"/>

	<space-output
		v-if="space"
		:space="space"
		:show-path="includeParentPath">

		<slot/>

	</space-output>

	<el-alert
		v-else title="Space could not be loaded"
		type="error"
		:closable="false"
		/>

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
		spaceId: {
			type: [String, Number],
			required: true,
		},
		includeParentPath: {
			type: Boolean,
			default: false,
		},
	},
	data() {
		return {
			loading: true,
			space: null,
		};
	},
	mounted() {
		ajaxGet('/ajax/space', {
			spaceId: this.spaceId,
			includeParentPath: this.includeParentPath,
		}).then(response => {
			this.space = response;
		}).finally(() => {
			this.loading = false;
		});
	},
};
</script>
