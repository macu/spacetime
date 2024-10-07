<template>
<div class="dashboard-page flex-column-lg page-width-md">

	<h2>Dashboard</h2>

	<loading-message v-if="loading"/>

	<category-link v-else-if="treetimeNode" :node="treetimeNode"/>

</div>
</template>

<script>
import CategoryLink from '@/widgets/category-link.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		CategoryLink,
	},
	data() {
		return {
			loading: true,
			treetimeNode: null,
		};
	},
	mounted() {
		this.loadDashboard();
	},
	methods: {
		loadDashboard() {
			this.loading = true;
			ajaxGet('/ajax/dashboard').then(response => {
				this.treetimeNode = response.treetimeNode;
			}).finally((error) => {
				this.loading = false;
			});
		},
	},
};
</script>
