<template>
<form-layout :title="root ? 'Create space' : 'Create branch'">

	<form-field :title="root ? 'Label for new space' : 'Label for new branch'" required>
		<template #tip>
			<div>A label is unique and permanent identifier for a given space. It cannot be changed later.</div>
		</template>
		<el-input
			v-model="label"
			:maxlength="$store.getters.labelMaxLength"
			show-word-limit
			size="large"
			:disabled="posting"
			/>
	</form-field>

	<form-field v-if="userAllowPin">
		<el-checkbox v-model="pin" :disabled="posting">
			Pin this space to the top of the parent space
		</el-checkbox>
	</form-field>

	<tags-field
		v-model:tags="tags"
		v-model:pinnedTags="pinnedTags"
		:allow-pinning="allowPinSubspaces"
		:disabled="posting"
	/>

	<form-actions>
		<el-button @click="submit()" type="primary" :disabled="createDisabled">
			Create
		</el-button>
	</form-actions>

</form-layout>
</template>

<script>
import TagsField from './tags-field.vue';

import {
	SPACE_TYPES,
} from '@/const.js';

export default {
	emits: ['submit'],
	components: {
		TagsField,
	},
	props: {
		root: {
			type: Boolean,
			default: false,
		},
		posting: {
			type: Boolean,
			default: false,
		},
		userAllowPin: {
			type: Boolean,
			default: false,
		},
		userAllowPinSubspacesOnCreate: {
			type: Boolean,
			default: false,
		},
	},
	data() {
		return {
			label: '',
			pin: false,
			tags: [],
			pinnedTags: [],
		};
	},
	computed: {
		createDisabled() {
			return this.posting || !this.label.trim();
		},
		allowPinSubspaces() {
			return this.userAllowPinSubspacesOnCreate &&
				this.userAllowPinSubspacesOnCreate(SPACE_TYPES.BRANCH);
		},
	},
	methods: {
		submit() {
			if (this.createDisabled) {
				return;
			}
			this.$emit('submit', {
				label: this.label.trim(),
				pin: this.allowUserSpaceActions && this.pin,
				tags: JSON.stringify(this.tags),
				pinnedTags: JSON.stringify(this.allowUserSpaceActions ? this.pinnedTags : []),
			});
		},
	},
};
</script>
