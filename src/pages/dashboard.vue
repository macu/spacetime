<template>
<div class="dashboard-page flex-column-lg page-width-md">

	<h2>Dashboard</h2>

	<loading-message v-if="loading"/>

	<node-header v-else-if="treetimeNode" :node="treetimeNode" link-to/>

</div>
</template>

<script>
import NodeHeader from '@/widgets/node-header.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		NodeHeader,
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
