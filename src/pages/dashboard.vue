<template>
<div class="dashboard-page flex-column-lg page-width-md">

	<el-backtop :right="40" :bottom="40" :visibility-height="1000"/>

	<create-dropdown :disabled="disableCreate" sticky/>

	<spaces-list :spaces="spaces" @load-more="loadMore()"/>

	<loading-message v-if="loading"/>

</div>
</template>

<script>
import CreateDropdown from '@/widgets/create-dropdown.vue';
import SpacesList from '@/widgets/spaces-list.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		CreateDropdown,
	},
	data() {
		return {
			loading: true,
			spaces: [],
		};
	},
	computed: {
		authenticated() {
			return this.$store.getters.userIsAuthenticated;
		},
		disableCreate() {
			return !this.authenticated;
		},
	},
	mounted() {
		this.loadDashboard();
	},
	methods: {
		loadDashboard() {
			this.loading = true;
			ajaxGet('/ajax/subspaces', {
				parentId: null, // root
			}).then(response => {
				this.spaces = response;
			}).finally((error) => {
				this.loading = false;
			});
		},
	},
};
</script>
