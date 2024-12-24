<template>
<div class="dashboard-page flex-column-lg page-width-md">

	<el-backtop :right="40" :bottom="40" :visibility-height="1000"/>

	<create-dropdown :disabled="$store.getters.createDisabled" sticky/>

	<spaces-list :spaces="spaces" :loading="loading" @load-more="loadMore()"/>

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
		SpacesList,
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
				limit: this.$store.getters.maxLimit,
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
				limit: this.$store.getters.maxLimit,
			}).then(response => {
				this.spaces = this.spaces.concat(response);
			}).finally((error) => {
				this.loading = false;
			});
		},
	},
};
</script>
