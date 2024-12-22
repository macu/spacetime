<template>
<div class="space-page page-width-lg">
	<loading-message v-if="loading"/>
	<space-output v-else-if="space" :space="space">
		<template #subspaces>
			<div @click.stop class="subspaces flex-column">
				<create-dropdown
					:parent-id="space.id"
					/>
				<space-list
					v-if="subspaces"
					:spaces="subspaces"
					:loading="loadingSubspaces"
					@load-more="loadMore()"
					/>
			</div>
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
import CreateDropdown from '@/widgets/create-dropdown.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		SpaceOutput,
		SpaceList,
		CreateDropdown,
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
			vm.loadSpace(to.params.spaceId);
		});
	},
	beforeRouteUpdate(to, from, next) {
		this.space = null;
		this.loading = true;
		next();
		this.loadSpace(to.params.spaceId);
	},
	methods: {
		loadSpace(spaceId) {
			this.loading = true;
			this.space = null;
			this.subspaces = null;
			ajaxGet('/ajax/space', {
				spaceId,
				includeSubspaces: true,
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

<style lang="scss">
.space-page {
	.subspaces {
		padding: 20px;
		background-color: darkturquoise;
		border-radius: 12px;
		cursor: default;
	}
}
</style>
