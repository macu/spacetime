<template>
<div class="category-view-page flex-column-lg page-width-md">

	<loading-message v-if="loadingHeader"/>

	<template v-else-if="header">

		<div v-if="path.length" class="field">
			<label>Path</label>
			<parent-path :path="path"/>
		</div>

		<h2>Category</h2>

		<category-header :node="header"/>

		<div class="flex-row-md">
			<el-button @click="addSubcategory()" type="primary">Add Subcategory</el-button>
			<el-button @click="addPost()" type="primary">Add post</el-button>
		</div>

		<loading-message v-if="loadingChildren"/>

		<nodes-list v-else :parent-id="header.id" :nodes="children"/>

	</template>

	<el-alert v-else
		title="Category could not be loaded."
		type="error" show-icon :closable="false"
		/>

</div>
</template>

<script>
import ParentPath from '@/widgets/parent-path.vue';
import CategoryHeader from '@/widgets/category-header.vue';
import NodesList from '@/widgets/nodes-list.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		ParentPath,
		CategoryHeader,
		NodesList,
	},
	data() {
		return {
			loadingHeader: true,
			header: null,
			path: [],

			loadingChildren: true,
			children: [],
		};
	},
	computed: {
		id() {
			return this.$route.params.id;
		},
	},
	beforeRouteEnter(to, from, next) {
		next(vm => {
			vm.loadCategory(to.params.id);
		});
	},
	beforeRouteUpdate(to, from, next) {
		this.loadingHeader = true;
		this.loadingChildren = true;
		this.header = null;
		this.path = [];
		this.children = [];

		next();

		this.loadCategory(to.params.id);
	},
	methods: {
		loadCategory(id) {
			this.loadingHeader = true;
			ajaxGet('/ajax/node/view', {
				id,
			}).then(response => {
				this.header = response.header;
				this.path = response.path;
				this.loadChildren(id);
			}).finally(() => {
				this.loadingHeader = false;
			});
		},
		loadChildren(id) {
			this.loadingChildren = true;
			ajaxGet('/ajax/node/children', {
				id,
			}).then(response => {
				this.children = response;
			}).finally(() => {
				this.loadingChildren = false;
			});
		},
	},
};
</script>
