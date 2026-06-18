<template>

	<loading-message v-if="loading"/>

	<space-output
		v-if="space"
		:space="space"
		:show-path="includeParentPath"
		:user-allow-pin="userAllowPin"
		:user-allow-pin-subspaces-on-create="userAllowPinSubspacesOnCreate"
		>

		<slot
			:space="space"
			:user-allow-pin="userAllowPin"
			:user-allow-pin-subspaces-on-create="userAllowPinSubspacesOnCreate"
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
		userAllowPin() {
			if (!this.space.parentPath.length) {
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
		// userAllowPinSubspaces() {
		// 	if (this.userAllowPin) {
		// 		return true;
		// 	}

		// 	// Allow pinning subspaces under root user space
		// 	let path = [...(this.space.parentPath || []), this.space];
		// 	let root = path[0];
		// 	if (root && root.spaceType === SPACE_TYPES.USER) {
		// 		return true;
		// 	}

		// 	return false;
		// },
		userAllowPinSubspacesOnCreate() {
			// Return a function that checks if pins are allowed under the given create space type
			return (createSpaceType) => {
				// If user space at root, allow pinning by user
				let path = [...(this.space.parentPath || []), this.space];
				let root = path[0];
				if (root && root.spaceType === SPACE_TYPES.USER) {
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
			}).finally(() => {
				this.loading = false;
			});
		},
	},
};
</script>
