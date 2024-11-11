<template>
<div class="dashboard-page flex-column-lg page-width-md">

	<h2>Dashboard</h2>

	<div class="flex-row">
		<create-dropdown :disabled="disableCreate"/>
		<el-button @click="gotoLanguages()">
			<material-icon icon="language"/>
			<span>Languages</span>
		</el-button>
		<el-button @click="gotoTags()">
			<material-icon icon="label"/>
			<span>Tags</span>
		</el-button>
	</div>

	<loading-message v-if="loading"/>

	<template v-else>

		<node-header v-if="treetimeNode" :node="treetimeNode" link-to/>

	</template>

</div>
</template>

<script>
import NodeHeader from '@/widgets/node-header.vue';
import CreateDropdown from '@/widgets/create-dropdown.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		NodeHeader,
		CreateDropdown,
	},
	data() {
		return {
			loading: true,
			treetimeNode: null,
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
			ajaxGet('/ajax/dashboard').then(response => {
				this.treetimeNode = response.treetimeNode;
			}).finally((error) => {
				this.loading = false;
			});
		},
		gotoCreateCategory() {
			this.$router.push({name: 'create-category'});
		},
		gotoLanguages() {
			this.$router.push({name: 'langs-view'});
		},
		gotoTags() {
			this.$router.push({name: 'tags-view'});
		},
	},
};
</script>
