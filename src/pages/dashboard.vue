<template>
<div class="dashboard-page flex-column-lg page-width-md">

	<h2>Dashboard</h2>

	<div class="flex-row">
		<el-button @click="gotoCreateCategory()" type="primary">
			<material-icon icon="add"/>
			<span>Create category</span>
		</el-button>
		<el-button @click="gotoLanguages()">
			<material-icon icon="language"/>
			<span>Languages</span>
		</el-button>
		<el-button @click="gotoTypes()">
			<material-icon icon="category"/>
			<span>Types</span>
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
		gotoCreateCategory() {
			this.$router.push({name: 'create-category'});
		},
		gotoLanguages() {
		},
		gotoTypes() {
		},
		gotoTags() {
		},
	},
};
</script>
