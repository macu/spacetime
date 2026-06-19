<template>

	<loading-message v-if="loading"/>

	<space-output
		v-else-if="space"
		:space="space"
		:show-path="includeParentPath"
		:permissions="permissions"
		@set-pinned="pinned => space.isPinned = pinned"
		>

		<slot
			:space="space"
			:permissions="permissions"
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
			required: true,
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

			if (!this.space || !this.space.parentPath) {
				return false;
			}

			// If user space at root, allow pinning by user
			let path = [...(this.space.parentPath || []), this.space];
			let root = path[0];
			if (root && root.spaceType === SPACE_TYPES.USER) {
				return root.createdBy === this.$store.getters.currentUserId;
			}

			// Otherwise, only allow pinning under text space by current user
			let textSpace = path.reverse().find(p => p.spaceType === SPACE_TYPES.TEXT);
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
			let textSpace = path.reverse().find(p => p.spaceType === SPACE_TYPES.TEXT);
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
				let textSpace = path.reverse().find(p => p.spaceType === SPACE_TYPES.TEXT);
				if (textSpace) {
					return textSpace.createdBy === this.$store.getters.currentUserId;
				}

				return false;
			};
		},
		permissions() {
			return {
				userAllowPinToParent: this.userAllowPinToParent,
				userAllowPinSubs: this.userAllowPinSubs,
				userAllowPinSubsOnCreate: this.userAllowPinSubsOnCreate,
			};
		},
	},
	watch: {
		spaceId() {
			this.loadSpace();
		},
	},
	mounted() {
		this.loadSpace();
	},
	methods: {
		loadSpace() {
			ajaxGet('/ajax/space', {
				spaceId: this.spaceId,
				includeParentPath: this.includeParentPath,
				includeTags: this.includeTags,
			}).then(response => {
				this.space = response;
				this.$emit('space-loaded', this.space);
			}).finally(() => {
				this.loading = false;
			});
		},
	},
};
</script>
