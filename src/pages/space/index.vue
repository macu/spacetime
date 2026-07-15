<template>
<div class="space-page page-width-xl">

	<return-to-top/>

	<space-loader
		:space-id="spaceId"
		:filter="filter"
		include-parent-path
		include-tags
		@space-loaded="handleSpaceLoaded">

		<template  #actions-area="{space, context}">
			<!-- tags output; tags are included in subspaces in pinned view -->
			<div class="flex-row-md">

				<el-button
					v-if="context.userAllowPinToParent"
					:type="space.isPinned ? 'success' : 'primary'"
					@click="togglePinned()">
					<material-icon v-if="space.isPinned" icon="keep"/>
					<material-icon v-else icon="keep_off"/>
					{{ space.isPinned ? 'Unpin' : 'Pin' }}
				</el-button>

				<checkin-button :space="space"/>

				<bookmark-button :space="space"/>

				<template v-if="!showingPinned">
					<add-tag-button
						:parent-id="spaceId"
						@added="tag => tags.push(tag)"
					/>
					<space-tag
						v-for="t in uniqueTags"
						:key="t.id"
						:space="t"
						:show-pinning="context && context.userAllowPinSubs"
						@toggle-pinned="toggleTagPinned(t)"
					/>
					<loading-message v-if="loadingTags" message="Loading tags..."/>
					<el-button v-else-if="showLoadMoreTags"
						@click="loadTags(true)" type="primary">
						Load more tags
					</el-button>
				</template>

			</div>
		</template>

		<template #default="{context}">

			<loading-message v-if="loadingSubspaces" message="Loading subspaces..."/>

			<div v-else @click.stop class="flex-column-lg">

				<horizontal-controls>
					<create-dropdown
						:parent-id="spaceId"
						:disabled="$store.getters.createDisabled"
						/>
					<subspaces-filter
						v-model="filter"
						allow-pinned
						/>
				</horizontal-controls>

				<div ref="subspaces" class="subspaces flex-column-lg">
					<subspace
						v-for="s in uniqueSubspaces"
						:key="s.id"
						:subspace="s"
						:context="context"
						:show-reorder="showReorderSubs"
						:filter-json="filterJson"
						@toggle-pinned="togglePinned(s)"
					/>
					<el-button v-if="showLoadMoreSubspaces"
						@click="loadSubspaces(true)" type="primary">
						Load more
					</el-button>
				</div>

				<el-alert v-if="showPinnedNotSupported" type="warning" :closable="false">
					<p>Pinned subspaces are not supported in this space.</p>
				</el-alert>

			</div>

		</template>

	</space-loader>

</div>
</template>

<script>
import SpaceLoader from '@/widgets/space-loader.vue';
import SpaceOutput from '@/widgets/space-output.vue';
import SpaceTag from '@/widgets/space-tag.vue';
import CreateDropdown from '@/widgets/create-dropdown.vue';
import CheckinButton from '@/widgets/checkin-button.vue';
import BookmarkButton from '@/widgets/bookmark-button.vue';
import AddTagButton from '@/widgets/add-tag-button.vue';
import Subspace from './subspace.vue';
import SubspacesFilter, {
	getFilter,
	FILTER_MODES,
} from '@/widgets/subspaces-filter.vue';

import {
	ajaxGet,
	ajaxPost,
} from '@/utils/ajax.js';

import {
	SPACE_TYPES,
} from '@/const.js';

import {
	alertSuccess,
} from '@/utils/notify.js';

import {
	loadScript,
	SCRIPTS,
} from '@/utils/loader.js';

import {
	startAutoscroll,
} from '@/utils/scroll.js';

export default {
	components: {
		SpaceLoader,
		SpaceOutput,
		SpaceTag,
		CreateDropdown,
		SubspacesFilter,
		CheckinButton,
		BookmarkButton,
		AddTagButton,
		Subspace,
	},
	setup() {
		return {
			// non-reactive properties
			drake: null,
			autoscroll: null,
		};
	},
	data() {
		return {
			filter: getFilter(),

			space: null,
			context: null,
			tags: [],
			subspaces: [],

			dragging: false,

			loadingTags: false,
			showLoadMoreTags: false,
			loadingSubspaces: false,
			showLoadMoreSubspaces: false,
		};
	},
	computed: {
		spaceId() {
			return this.$route.params.spaceId;
		},
		filterJson() {
			return this.filter ? JSON.stringify(this.filter) : null;
		},
		showingPinned() {
			return this.filter &&
				this.filter.mode === FILTER_MODES.PINNED;
		},
		showPinnedNotSupported() {
			return this.showingPinned &&
				this.context &&
				!this.context.supportsPinning;
		},
		showReorderSubs() {
			return this.showingPinned &&
				this.context &&
				this.context.userAllowPinSubs &&
				!this.loadingSubspaces;
		},
		uniqueTags() {
			const seen = new Set();
			return this.tags.filter(t => {
				if (seen.has(t.id)) {
					return false;
				}
				seen.add(t.id);
				return true;
			});
		},
		uniqueSubspaces() {
			const seen = new Set();
			return this.subspaces.filter(s => {
				if (seen.has(s.id)) {
					return false;
				}
				seen.add(s.id);
				return true;
			});
		},
	},
	watch: {
		filter: {
			deep: true,
			handler() {
				this.reloadSubs();
			},
		},
		showReorderSubs: {
			immediate: true,
			handler(show) {
				if (show) {
					this.$nextTick(() => {
						this.loadDragula();
					});
				}
			},
		},
		dragging(dragging) {
			if (dragging) {
				if (!this.autoscroll) {
					this.autoscroll = startAutoscroll();
				}
			} else {
				if (this.autoscroll) {
					this.autoscroll.stop();
					this.autoscroll = null;
				}
			}
		},
	},
	beforeUnmount() {
		if (this.drake) {
			this.drake.destroy();
			this.drake = null;
		}
		if (this.autoscroll) {
			this.autoscroll.stop();
			this.autoscroll = null;
		}
	},
	methods: {

		// called from space-loader
		handleSpaceLoaded({space, context}) {
			this.space = space;
			this.context = context;
			this.tags = space.tags;
			this.subspaces = space.subspaces;
			this.showLoadMoreTags = space.tags.length === this.$const.defaultTagsLimit;
			this.showLoadMoreSubspaces = space.subspaces.length === this.$const.maxPageLimit;
		},

		// filters changed
		reloadSubs() {
			this.loadingTags = true;
			this.loadingSubspaces = true;
			ajaxGet('/ajax/space/reload', {
				spaceId: this.spaceId,
				filter: this.filterJson,
			}).then(response => {
				this.tags = response.tags || [];
				this.subspaces = response.subspaces || [];
			}).finally(() => {
				this.loadingTags = false;
				this.loadingSubspaces = false;
			});
		},

		// loading more tags
		loadTags(more = false) {
			this.loadingTags = true;
			if (!more) {
				this.tags = [];
			}
			ajaxGet('/ajax/space/tags', {
				parentId: this.spaceId,
				offset: more ? this.tags.length : 0,
				limit: this.$const.defaultTagsLimit,
				filter: this.filterJson,
			}).then(response => {
				if (more) {
					this.tags = this.tags.concat(response);
				} else {
					this.tags = response;
				}
				this.showLoadMoreTags = response.length === this.$const.defaultTagsLimit;
			}).finally(() => {
				this.loadingTags = false;
			});
		},

		// loading more subspaces
		loadSubspaces(more = false) {
			this.loadingSubspaces = true;
			if (!more) {
				this.subspaces = [];
			}
			ajaxGet('/ajax/subspaces', {
				parentId: this.spaceId,
				offset: more ? this.subspaces.length : 0,
				limit: this.$const.maxPageLimit,
				filter: this.filterJson,
			}).then(response => {
				this.subspaces = this.subspaces.concat(response);
			}).finally(() => {
				this.loadingSubspaces = false;
			});
		},

		togglePinned(space = null) {
			if (!space) {
				space = this.space;
			} else {
				space = this.subspaces.find(s => s.id === space.id);
			}
			let pinned = !space.isPinned;
			if (!pinned) {
				// Always confirm - unpinning causes a pinned space to lose its position
				this.$confirm('Are you sure you want to unpin this space?', 'Confirm unpin', {
					confirmButtonText: 'Yes',
					cancelButtonText: 'No',
					type: 'warning',
				}).then(() => {
					space.isPinned = pinned;
					ajaxPost('/ajax/space/pin', {
						spaceId: space.id,
						pinned,
					}).then(response => {
						alertSuccess('Space unpinned');
					}).catch(err => {
						space.isPinned = !pinned; // revert
					});
				}).catch(() => {
					// Cancelled
				});
			} else {
				space.isPinned = pinned;
				ajaxPost('/ajax/space/pin', {
					spaceId: space.id,
					pinned,
				}).then(response => {
					alertSuccess('Space pinned');
				}).catch(err => {
					space.isPinned = !pinned; // revert
				});
			}
		},

		toggleTagPinned(t) {
			let pinned = !t.isPinned;

			if (!pinned) {
				// Always confirm - unpinning causes a pinned space to lose its position
				this.$confirm('Are you sure you want to unpin this tag?', 'Confirm unpin', {
					confirmButtonText: 'Yes',
					cancelButtonText: 'No',
					type: 'warning',
				}).then(() => {
					t.isPinned = pinned;
					ajaxPost('/ajax/space/pin', {
						spaceId: t.id,
						pinned,
					}).then(() => {
						this.$message.success('Tag unpinned');
					}).catch(err => {
						t.isPinned = !pinned; // revert
						this.$message.error('Failed to unpin tag');
					});
				}).catch(() => {
					// Cancelled
				});
				return;
			} else {
				t.isPinned = pinned;
				ajaxPost('/ajax/space/pin', {
					spaceId: t.id,
					pinned,
				}).then(() => {
					this.$message.success('Tag pinned');
				}).catch(err => {
					t.isPinned = !pinned; // revert
					this.$message.error('Failed to pin tag');
				});
			}
		},

		loadDragula() {
			loadScript(SCRIPTS.DRAGULA).then(Dragula => {
				this.drake = Dragula([this.$refs.subspaces], {
					moves: (el, source, handle, sibling) => {
						return handle.classList.contains('drag-handle');
					},
				}).on('drag', () => {
					this.dragging = true;
				}).on('dragend', () => {
					this.dragging = false;
				}).on('drop', (el, target, source, sibling) => {
					let spaceId = el.getAttribute('data-space-id');
					let beforeId = sibling ? sibling.getAttribute('data-space-id') : null;

					this.drake.cancel(true); // revert, reorder in data

					let existingIndex = this.subspaces.findIndex(s => (''+s.id) === spaceId);

					let pinnedSpace = this.subspaces.splice(existingIndex, 1)[0];

					let newIndex = sibling
						? this.subspaces.findIndex(s => (''+s.id) === beforeId)
						: this.subspaces.length;

					console.log('Move from', existingIndex, 'to', newIndex);

					// Move in data
					this.subspaces.splice(newIndex, 0, pinnedSpace);

					ajaxPost('/ajax/space/move-pin', {
						spaceId,
						newIndex,
					}).then(() => {
						this.$message.success('Pinned space moved');
					}).catch(err => {
						// Revert on error
						this.subspaces.splice(newIndex, 1);
						this.subspaces.splice(existingIndex, 0, pinnedSpace);
						this.$message.error('Failed to move pinned space');
					});
				});
			}).catch(err => {
				this.$message.error('Failed to load drag and drop functionality');
			});
		},

	},
};
</script>
