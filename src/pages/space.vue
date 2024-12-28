<template>
<div class="space-page page-width-lg">

	<loading-message v-if="loading"/>

	<space-output v-else-if="space" :space="space" show-path>

		<div @click.stop class="subspaces flex-column-md">

			<create-dropdown
				:parent-id="space.id"
				:disabled="$store.getters.createDisabled"
				/>

			<space-list
				v-if="subspaces"
				:spaces="subspaces"
				:loading="loadingSubspaces"
				@load-more="loadMore()"
				/>

		</div>

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

import {
	SPACE_TYPES,
} from '@/const.js';

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
	watch: {
		spaceId: {
			immediate: true,
			handler() {
				this.$nextTick(() => {
					this.loadSpace(this.spaceId);
				});
			},
		},
	},
	methods: {
		loadSpace(spaceId) {
			this.loading = true;
			this.space = null;
			this.subspaces = null;
			ajaxGet('/ajax/space', {
				spaceId,
				includeTags: true,
				includeSubspaces: true,
				includeParentPath: true,
				excludeTypes: SPACE_TYPES.CHECK_IN,
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
				excludeTypes: SPACE_TYPES.CHECK_IN,
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
		// padding: 20px;
		// background-color: darkturquoise;
		// border-radius: 12px;
		// cursor: default;
	}
}
</style>
