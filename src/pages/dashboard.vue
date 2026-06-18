<template>
<div class="dashboard-page flex-column-lg page-width-md">

	<return-to-top/>

	<create-dropdown :disabled="$store.getters.createDisabled" sticky/>

	<div class="flex-column-lg">

		<space-output
			v-for="s in spaces"
			:space="s"
		/>

		<div class="center">
			<loading-message v-if="loading"/>
			<el-button v-else @click="loadMore()" type="primary">Load More</el-button>
		</div>

	</div>

</div>
</template>

<script>
import CreateDropdown from '@/widgets/create-dropdown.vue';
import SpaceOutput from '@/widgets/space-output.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		CreateDropdown,
		SpaceOutput,
	},
	data() {
		return {
			loading: true,
			spaces: [],
		};
	},
	mounted() {
		this.loadDashboard();
	},
	methods: {
		loadDashboard() {
			this.loading = true;
			ajaxGet('/ajax/subspaces', {
				parentId: null, // root
				offset: 0,
				limit: this.$store.getters.maxPageLimit,
				includeTags: true,
			}).then(response => {
				this.spaces = response;
			}).finally((error) => {
				this.loading = false;
			});
		},
		loadMore() {
			this.loading = true;
			ajaxGet('/ajax/subspaces', {
				parentId: null, // root
				offset: this.spaces.length,
				limit: this.$store.getters.maxPageLimit,
				includeTags: true,
			}).then(response => {
				this.spaces = this.spaces.concat(response);
			}).finally((error) => {
				this.loading = false;
			});
		},
	},
};
</script>
