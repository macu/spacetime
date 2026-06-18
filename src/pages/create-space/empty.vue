<template>
<div class="create-branch-space-page flex-column-lg page-width-md">

	<space-loader v-if="parentId" :space-id="parentId" include-parent-path>

		<template #default="{userAllowPin, userAllowPinSubspacesOnCreate}">

			<form-fields
				:posting="posting"
				:root="root"
				:user-allow-pin="userAllowPin"
				:user-allow-pin-subspaces-on-create="userAllowPinSubspacesOnCreate"
				@submit="submit"
			/>

		</template>

	</space-loader>

	<form-fields v-else :posting="posting" :root="root" @submit="submit"/>

</div>
</template>

<script>
import SpaceLoader from '@/widgets/space-loader.vue';
import FormFields from './empty-form.vue';

import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	components: {
		SpaceLoader,
		FormFields,
	},
	data() {
		return {
			posting: false,
		};
	},
	computed: {
		parentId() {
			return this.$route.query.parentId ?
				parseInt(this.$route.query.parentId, 10) : null;
		},
		root() {
			return this.parentId === null;
		},
	},
	methods: {
		submit(payload) {
			this.posting = true;
			ajaxPost('/ajax/space/create/branch', {
				parentId: this.parentId,
				...payload,
			}, {
				409: 'The given label already exists at this level.',
			}).then(response => {
				this.$router.replace({
					name: 'space',
					params: {
						spaceId: response.id,
					},
				});
			}).catch(() => {
				this.posting = false;
			});
		},
	},
};
</script>
