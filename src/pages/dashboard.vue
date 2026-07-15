<template>
<div class="dashboard-page flex-column-lg page-width-md">

	<return-to-top/>

	<horizontal-controls>
		<create-dropdown :disabled="$store.getters.createDisabled" sticky/>
		<subspaces-filter
			v-model="filter"
			:allow-pinned="false"
		/>
	</horizontal-controls>

	<div class="flex-column-lg">

		<space-output
			v-for="s in spaces"
			:space="s">
			<template #actions-area>
				<div class="flex-row-md">
					<checkin-button :space="s"/>
					<bookmark-button :space="s"/>
					<space-tag
						v-for="t in s.tags"
						:key="t.id"
						:space="t"
						:show-pinning="false"
					/>
				</div>
			</template>
		</space-output>

		<div class="center">
			<loading-message v-if="loading"/>
			<el-button v-else @click="loadMore()" type="primary">Load More</el-button>
		</div>

	</div>

</div>
</template>

<script>
import CreateDropdown from '@/widgets/create-dropdown.vue';
import SpaceOutput from '@/widgets/space-output.vue';
import CheckinButton from '@/widgets/checkin-button.vue';
import BookmarkButton from '@/widgets/bookmark-button.vue';
import SpaceTag from '@/widgets/space-tag.vue';
import SubspacesFilter, {
	getFilter,
} from '@/widgets/subspaces-filter.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		CreateDropdown,
		SpaceOutput,
		CheckinButton,
		BookmarkButton,
		SpaceTag,
		SubspacesFilter,
	},
	data() {
		return {
			loading: true,
			spaces: [],
			filter: getFilter(false),
		};
	},
	watch: {
		filter: {
			handler() {
				this.loadDashboard();
			},
			deep: true,
		},
	},
	mounted() {
		this.loadDashboard();
	},
	methods: {
		loadDashboard() {
			this.loading = true;
			this.spaces = [];
			ajaxGet('/ajax/subspaces', {
				parentId: null, // root
				offset: 0,
				limit: this.$const.maxPageLimit,
				includeTags: true,
				filter: this.filter ? JSON.stringify(this.filter) : null,
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
				limit: this.$const.maxPageLimit,
				includeTags: true,
				filter: this.filter ? JSON.stringify(this.filter) : null,
			}).then(response => {
				this.spaces = this.spaces.concat(response);
			}).finally((error) => {
				this.loading = false;
			});
		},
	},
};
</script>
