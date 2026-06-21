<template>
<div class="space-page page-width-xl">

	<return-to-top/>

	<space-loader
		:space-id="spaceId"
		include-parent-path include-tags
		@space-loaded="data => handleSpaceLoaded(data)">

		<template #default="{permissions}">

			<loading-message v-if="loadingSubspaces" message="Loading subspaces..."/>

			<div v-else @click.stop class="flex-column-lg">

				<div class="flex-row-md">
					<create-dropdown
						:parent-id="spaceId"
						:disabled="$store.getters.createDisabled"
						/>
					<subspaces-filter
						v-model="filter"
						allow-pinned
						/>
				</div>

				<div ref="subspaces" class="subspaces flex-column-lg">
					<space-output
						v-for="s in subspaces"
						:space="s"
						:permissions="permissions"
						sub-space
						:show-reorder="showReorderSubs"
						:data-space-id="s.id"
						@set-pinned="pinned => setPinned(s, pinned)"
					/>
				</div>

			</div>

		</template>

	</space-loader>

</div>
</template>

<script>
import SpaceLoader from '@/widgets/space-loader.vue';
import SpaceOutput from '@/widgets/space-output.vue';
import CreateDropdown from '@/widgets/create-dropdown.vue';
import SubspacesFilter, {getFilter, FILTER_MODES} from '@/widgets/subspaces-filter.vue';

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
		CreateDropdown,
		SubspacesFilter,
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
			loadingSubspaces: true,
			subspaces: [],
			permissions: null,
			dragging: false,
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
			return this.filter && this.filter.mode === FILTER_MODES.PINNED;
		},
		showReorderSubs() {
			return this.showingPinned &&
				this.permissions &&
				this.permissions.userAllowPinSubs &&
				!this.loadingSubspaces;
		},
	},
	watch: {
		filter: {
			deep: true,
			handler() {
				this.loadSubspaces();
			},
		},
		spaceId: {
			immediate: true,
			handler() {
				this.loadSubspaces();
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
		handleSpaceLoaded({space, permissions}) {
			this.permissions = permissions;
		},

		loadSubspaces(more = false) {
			this.loadingSubspaces = true;

			if (!more) {
				this.subspaces = [];
			}

			ajaxGet('/ajax/subspaces', {
				parentId: this.spaceId,
				offset: this.subspaces.length,
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
