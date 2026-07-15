<template>

	<loading-message v-if="loading"/>

	<space-output
		v-else-if="space"
		:space="space"
		:show-path="includeParentPath"
		:context="context"
		@set-pinned="pinned => space.isPinned = pinned"
		>

		<template #actions-area="{space, context}">
			<slot name="actions-area" :space="space" :context="context"/>
		</template>

		<slot
			:space="space"
			:context="context"
		/>

	</space-output>

	<el-alert
		v-else title="Space could not be loaded"
		type="error"
		:closable="false"
		/>

</template>

<script>
import SpaceOutput from '@/widgets/space-output.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

import {
	SPACE_TYPES,
	USER_CONTEXT_TYPES,
} from '@/const.js';

export default {
	emits: [
		'space-loaded',
	],
	components: {
		SpaceOutput,
	},
	props: {
		spaceId: {
			type: [String, Number],
			required: false,
		},
		filter: {
			type: Object,
			required: false,
		},
		includeParentPath: {
			type: Boolean,
			default: false,
		},
		includeTags: {
			type: Boolean,
			default: false,
		},
	},
	data() {
		return {
			loading: true,
			space: null,
		};
	},
	computed: {
		userAllowPinToParent() {
			// Allow pinning loaded space to parent

			if (!this.space || !this.space.parentPath || this.space.parentPath.length === 0) {
				return false;
			}

			// If user space at root, allow pinning by user
			let path = [...(this.space.parentPath || []), this.space];
			let root = path[0];
			if (root && root.spaceType === SPACE_TYPES.USER) {
				return root.createdBy === this.$store.getters.currentUserId;
			}

			// Otherwise, only allow pinning under text space by current user
			let textSpace = [...(this.space.parentPath || [])].reverse().find(
				p => p.spaceType === SPACE_TYPES.TEXT
			);
			if (textSpace) {
				return textSpace.createdBy === this.$store.getters.currentUserId;
			}

			return false;
		},
		userAllowPinSubs() {
			// Allow pinning subspaces within loaded space

			if (!this.space) {
				return false;
			}

			// If user space at root, allow pinning by user
			let path = [...(this.space.parentPath || []), this.space];
			let root = path[0];
			if (root.spaceType === SPACE_TYPES.USER) {
				return root.createdBy === this.$store.getters.currentUserId;
			}

			// In public spaces, only allow pinning within text space by current user
			let textSpace = [...path].reverse().find(p => p.spaceType === SPACE_TYPES.TEXT);
			if (textSpace) {
				return textSpace.createdBy === this.$store.getters.currentUserId;
			}

			return false;
		},
		userAllowPinSubsOnCreate() {
			// Allow pinning subspaces to subspace being created under loaded space

			return (createSpaceType) => {
				if (!this.space) {
					return false;
				}

				// If user space at root, allow pinning by user
				let path = [...(this.space.parentPath || []), this.space];
				let root = path[0];
				if (root.spaceType === SPACE_TYPES.USER) {
					return root.createdBy === this.$store.getters.currentUserId;
				}

				// Always allow authors to pin under public text spaces
				if (createSpaceType === SPACE_TYPES.TEXT) {
					return true;
				}

				// Otherwise, only allow pinning within text space by current user
				let textSpace = [...path].reverse().find(p => p.spaceType === SPACE_TYPES.TEXT);
				if (textSpace) {
					return textSpace.createdBy === this.$store.getters.currentUserId;
				}

				return false;
			};
		},
		supportsPinning() {
			// return whether in user space or under a text space

			if (!this.space) {
				return false;
			}

			// If user space at root, space supports pinning
			let path = [...(this.space.parentPath || []), this.space];
			let root = path[0];
			if (root && root.spaceType === SPACE_TYPES.USER) {
				return true;
			}

			// Otherwise, space under a text space supports pinning
			let textSpace = [...(this.space.parentPath || [])].reverse().find(
				p => p.spaceType === SPACE_TYPES.TEXT
			);
			if (textSpace) {
				return true;
			}

			return false;
		},
		context() {
			return {
				userAllowPinToParent: this.userAllowPinToParent,
				userAllowPinSubs: this.userAllowPinSubs,
				userAllowPinSubsOnCreate: this.userAllowPinSubsOnCreate,
				supportsPinning: this.supportsPinning,
			};
		},
	},
	watch: {
		spaceId: {
			immediate: true,
			handler() {
				this.loadSpace();
			},
		},
	},
	methods: {
		loadSpace() {
			this.loading = true;
			ajaxGet('/ajax/space', {
				spaceId: this.spaceId,
				includeParentPath: this.includeParentPath,
				includeTags: this.includeTags,
				filter: this.filter ? JSON.stringify(this.filter) : null,
			}).then(response => {
				this.space = response;
				this.$emit('space-loaded', {space: response, context: this.context});
			}).finally(() => {
				this.loading = false;
			});
		},

		updateSpace(space) {
			this.space = space;
		},
	},
};
</script>
