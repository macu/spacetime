<template>
<div class="space-page page-width-xl">

	<return-to-top/>

	<space-loader
		:space-id="spaceId"
		:filter="filter"
		include-parent-path
		include-tags
		@space-loaded="handleSpaceLoaded">

		<template v-if="!showingPinned" #tags-area="{context}">
			<!-- tags output; tags are included in subspaces in pinned view -->
			<div class="flex-row-md">
				<add-tag-button
					:parent-id="spaceId"
					@added="tag => tags.push(tag)"
				/>
				<space-tag
					v-for="t in uniqueTags"
					:key="t.id"
					:space="t"
					:show-pinning="context && context.userAllowPinSubs"
					@set-pinned="pinned => setTagPinned(t, pinned)"
				/>
				<loading-message v-if="loadingTags" message="Loading tags..."/>
				<el-button v-else-if="showLoadMoreTags"
					@click="loadTags(true)" type="primary">
					Load more
				</el-button>
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
					<space-output
						v-for="s in uniqueSubspaces"
						:key="s.id"
						:space="s"
						:context="context"
						sub-space
						:show-reorder="showReorderSubs"
						:data-space-id="s.id"
						@set-pinned="pinned => setPinned(s, pinned)">
						<template #tags-area>
							<div class="flex-row-md">
								<space-tag
									v-for="t in s.tags"
									:space="t"
									/>
							</div>
						</template>
					</space-output>
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
import AddTagButton from '@/widgets/add-tag-button.vue';
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
		AddTagButton,
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
		filterJSON() {
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
			this.context = context;
			this.tags = space.tags;
			this.subspaces = space.subspaces;
		},

		// filters changed
		reloadSubs() {
			this.loadingTags = true;
			this.loadingSubspaces = true;
			ajaxGet('/ajax/space/reload', {
				spaceId: this.spaceId,
				filter: this.filterJSON,
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
				limit: 10,
				filter: this.filterJSON,
			}).then(response => {
				if (more) {
					this.tags = this.tags.concat(response);
				} else {
					this.tags = response;
				}
				this.showLoadMoreTags = response.length === 10;
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
				limit: this.$store.getters.maxPageLimit,
				filter: this.filterJSON,
			}).then(response => {
				this.subspaces = this.subspaces.concat(response);
			}).finally(() => {
				this.loadingSubspaces = false;
			});
		},

		setPinned(s, pinned) {
			s.isPinned = pinned;

			if (this.showingPinned && !pinned) {
				// If unpinning while showing pinned, remove from list
				this.subspaces = this.subspaces.filter(sub => sub.id !== s.id);
			}
		},

		setTagPinned(t, pinned) {
			t.isPinned = pinned;

			ajaxPost('/ajax/space/pin', {
				spaceId: t.id,
				pinned,
			}).then(() => {
				this.$message.success('Tag ' + (pinned ? 'pinned' : 'unpinned'));
			}).catch(err => {
				t.isPinned = !pinned; // revert
				this.$message.error('Failed to ' + (pinned ? 'pin' : 'unpin') + ' tag');
			});
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
