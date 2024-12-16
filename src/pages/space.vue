<template>
<div class="space-page page-width-md">
	<loading-message v-if="loading"/>
	<space-output v-else-if="space" :space="space">
		<template #subspaces>
			<space-list
				v-if="subspaces"
				:spaces="subspaces"
				:loading="loadingSubspaces"
				@load-more="loadMore()"
				/>
		</template>
	</space-output>
	<el-alert v-else type="error" show-icon :closable="false">
		This space could not be loaded.
	</el-alert>
</div>
</template>

<script>
import SpaceOutput from '@/widgets/space-output.vue';
import SpaceList from '@/widgets/spaces-list.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		SpaceOutput,
		SpaceList,
	},
	data() {
		return {
			loading: true,
			space: null,
			subspaces: null,
			loadingSubspaces: false,
		};
	},
	computed: {
		spaceId() {
			return this.$route.params.spaceId;
		},
	},
	beforeRouteEnter(to, from, next) {
		next(vm => {
			vm.loadSpace();
		});
	},
	beforeRouteUpdate(to, from, next) {
		this.space = null;
		this.loading = true;
		next();
		this.loadSpace();
	},
	methods: {
		loadSpace() {
			this.loading = true;
			this.space = null;
			this.subspaces = null;
			ajaxGet('/ajax/space', {
				spaceId: this.spaceId,
			}).then(response => {
				this.space = response;
				this.subspaces = response.topSubspaces
					? response.topSubspaces.slice() : null;
			}).finally((error) => {
				this.loading = false;
			});
		},
		loadMore() {
			if (!this.space) {
				return;
			}
			this.loadingSubspaces = true;
			ajaxGet('/ajax/subspaces', {
				parentId: this.space.id,
				offset: this.subspaces.length,
				limit: this.$store.getters.maxLimit,
			}).then(response => {
				this.subspaces = this.subspaces.concat(response);
			}).finally(() => {
				this.loadingSubspaces = false;
			});
		},
	},
};
</script>
