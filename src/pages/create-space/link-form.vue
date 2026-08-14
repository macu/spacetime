<template>
<form-layout title="Create link" class="create-link-form">

	<p>Select a space from your bookmarks to add as a link in this space.</p>

	<div v-if="selectedSpace" class="selected-space flex-column">
		<p><strong>Selected space:</strong></p>
		<space-output
			:space="selectedSpace"
			show-path
			:goto-space-on-click="false"
			:goto-path-on-click="false">
			<template #actions-area>
				<div>
					<el-button @click="selectedSpaceId = null" type="primary">
						Change selection
					</el-button>
				</div>
			</template>
		</space-output>
	</div>

	<loading-message v-else-if="loading" message="Loading bookmarks..."/>

	<div v-else-if="bookmarks.length > 0" class="bookmarks flex-column-lg">
		<space-output
			v-for="space in bookmarks"
			:key="space.id"
			:space="space"
			show-path
			:goto-space-on-click="false"
			:goto-path-on-click="false">
			<template #actions-area>
				<div class="flex-row">
					<el-alert v-if="space.parentId === parentId" type="warning" :closable="false">
						<p>A space cannot be linked into its own parent space.</p>
					</el-alert>
					<el-alert v-else-if="space.includedInParent" type="warning" :closable="false">
						<p>A link to this space has already been added in this parent space.</p>
					</el-alert>
					<el-alert v-else-if="space.spaceType === linkType" type="warning" :closable="false">
						<p>You cannot link a link.</p>
					</el-alert>
					<el-button v-else @click="selectedSpaceId = space.id" type="primary">
						Select
					</el-button>
				</div>
			</template>
		</space-output>
		<loading-message v-if="loadingMore" message="Loading more..."/>
		<el-button v-else-if="showLoadMore" @click="loadMore()" type="primary">
			Load more
		</el-button>
	</div>

	<el-alert v-else type="info" :closable="false">
		<p>You have no bookmarks yet. You can bookmark a space by clicking the bookmark icon in the space header.</p>
	</el-alert>

	<template v-if="bookmarks.length > 0">

		<!-- only allow pinning within author scope -->
		<form-field v-if="context.userAllowPinToParent">
			<el-checkbox v-model="pin" :disabled="posting">
				Pin this space to the top of the parent space
			</el-checkbox>
		</form-field>

		<!-- always allow author to pin under their own text space -->
		<tags-field
			v-model:tags="tags"
			v-model:pinnedTags="pinnedTags"
			:allow-pinning="allowPinSubspaces"
			:disabled="posting"
		/>

	</template>

	<form-actions>
		<el-button @click="submit()" type="primary" :disabled="createDisabled">
			Create
		</el-button>
	</form-actions>

</form-layout>
</template>

<script>
import TagsField from './tags-field.vue';
import SpaceOutput from '@/widgets/space-output.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

import {
	SPACE_TYPES,
} from '@/const.js';

export default {
	components: {
		TagsField,
		SpaceOutput,
	},
	props: {
		parentId: {
			type: Number,
			required: true,
		},
		posting: {
			type: Boolean,
			default: false,
		},
		context: {
			type: Object,
			required: true,
		},
	},
	data() {
		return {
			pin: false,
			tags: [],
			pinnedTags: [],

			loading: false,
			loadingMore: false,
			bookmarks: [],
			showLoadMore: false,
			selectedSpaceId: null,
		};
	},
	computed: {
		linkType() {
			return SPACE_TYPES.LINK;
		},
		createDisabled() {
			return this.posting || !this.selectedSpaceId;
		},
		allowPinSubspaces() {
			return this.context.userAllowPinSubsOnCreate(SPACE_TYPES.LINK);
		},
		selectedSpace() {
			if (!this.selectedSpaceId) {
				return null;
			}
			return this.bookmarks.find(s => s.id === this.selectedSpaceId) || null;
		},
	},
	mounted() {
		this.loadMore();
	},
	methods: {
		loadMore() {
			if (this.loading || this.loadingMore) {
				return;
			}

			if (this.bookmarks.length) {
				this.loadingMore = true;
			} else {
				this.loading = true;
			}

			ajaxGet('/ajax/bookmarks', {
				offset: this.bookmarks.length,
				limit: this.$const.maxPageLimit,
				includeParentPath: true,
				includeTags: true,
				includeLinkedInParentId: this.parentId,
			}).then(response => {
				this.bookmarks = this.bookmarks.concat(response);
				this.showLoadMore = response.length >= this.$const.maxPageLimit;
			}).finally(() => {
				this.loading = false;
				this.loadingMore = false;
			});
		},

		selectSpace(s) {
			if (s.includedInParent || s.spaceType === SPACE_TYPES.LINK) {
				// Can't link links
				return;
			}
			this.selectedSpaceId = s.id;
		},

		submit() {
			if (this.createDisabled) {
				return;
			}
			this.$emit('submit', {
				spaceId: this.selectedSpaceId,
				pin: this.context.userAllowPinToParent && this.pin,
				tags: JSON.stringify(this.tags),
				pinnedTags: JSON.stringify(this.pinnedTags), // pinned tags always allowd on text spaces
			});
		},
	},
};
</script>

<style lang="scss">
@import '@/styles/vars.scss';

.create-link-form {

	.selected-space {
		border: thin solid var(--el-color-primary);
		border-radius: $border-radius;
		padding: 1rem;
	}

	.bookmarks {
		border: thin solid var(--el-color-primary);
		border-radius: $border-radius;
		padding: 1rem;

		max-height: 80vh;
		overflow-y: auto;
	}

}
</style>
