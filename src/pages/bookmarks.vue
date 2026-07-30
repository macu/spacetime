<template>
<div class="bookmarks-page page-width-xl">

	<el-alert v-if="!authenticated"
		title="You must be logged in to view your bookmarks."
		type="warning"
		effect="dark"
		show-icon
		:closable="false"
	/>

	<loading-message v-else-if="loading" message="Loading bookmarks..."/>

	<div v-else class="flex-column-lg">

		<h1>Your bookmarks</h1>

		<template v-if="bookmarks.length > 0">
			<space-output
				v-for="space in bookmarksList"
				:key="space.id"
				:space="space"
				show-path>
				<template #actions-area>
					<div class="flex-row-md">
						<checkin-button :space="space"/>
						<bookmark-button
							:space="space"
							@bookmark-removed="removeBookmark(space)"
						/>
						<span v-text="space.bookmarkedAtText"/>
						<space-tag
							v-for="t in space.tags"
							:key="t.id"
							:space="t"
						/>
					</div>
				</template>
			</space-output>

			<loading-message v-if="loadingMore" message="Loading bookmarks..."/>

			<el-button v-else-if="showLoadMore" @click="loadMore" type="primary">
				Load more
			</el-button>
		</template>

		<el-alert v-else
			title="You have no bookmarks yet."
			type="info"
			effect="dark"
			show-icon
			:closable="false"
		/>

	</div>

</div>
</template>

<script>
import SpaceOutput from '@/widgets/space-output.vue';
import CheckinButton from '@/widgets/checkin-button.vue';
import BookmarkButton from '@/widgets/bookmark-button.vue';
import SpaceTag from '@/widgets/space-tag.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

import {
	alertSuccess,
} from '@/utils/notify.js';

export default {
	components: {
		SpaceOutput,
		CheckinButton,
		BookmarkButton,
		SpaceTag,
	},
	data() {
		return {
			loading: true,
			loadingMore: false,
			bookmarks: [],
			showLoadMore: false,
		};
	},
	computed: {
		authenticated() {
			return this.$store.getters.authenticated;
		},
		bookmarksList() {
			return this.bookmarks.map(space => {
				return {
					...space,
					bookmarkedAtText: space.bookmarkCreatedAt
						? "Bookmarked " + window.moment(space.bookmarkCreatedAt).fromNow()
						: null,
				};
			});
		},
	},
	watch: {
		authenticated(newVal, oldVal) {
			if (newVal && !oldVal) {
				this.loadBookmarks();
			}
		},
	},
	mounted() {
		if (this.authenticated) {
			this.loadBookmarks();
		}
	},
	methods: {
		loadBookmarks() {
			this.loading = true;
			ajaxGet('/ajax/bookmarks', {
				offset: 0,
				limit: this.$const.maxPageLimit,
				includeParentPath: true,
				includeTags: true,
			}).then(response => {
				this.bookmarks = response;
				this.showLoadMore = response.length >= this.$const.maxPageLimit;
			}).finally((error) => {
				this.loading = false;
			});
		},
		loadMore() {
			this.loadingMore = true;
			ajaxGet('/ajax/bookmarks', {
				offset: this.bookmarks.length,
				limit: this.$const.maxPageLimit,
				includeParentPath: true,
				includeTags: true,
			}).then(response => {
				this.bookmarks = this.bookmarks.concat(response);
				this.showLoadMore = response.length >= this.$const.maxPageLimit;
			}).finally((error) => {
				this.loadingMore = false;
			});
		},

		removeBookmark(space) {
			this.bookmarks = this.bookmarks.filter(s => s.id !== space.id);
			alertSuccess('Bookmark removed');
		},
	},
};
</script>

<style lang="scss">
@import '@/styles/vars.scss';

.bookmarks-page {
	background-color: $space-bg-color;
	border: thin solid darkblue;
	border-radius: $border-radius;
	padding: 40px;
	border-radius: 12px;
}
</style>
