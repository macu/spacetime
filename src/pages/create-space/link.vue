<template>
<div class="create-link-page flex-column-lg page-width-md">

	<el-alert v-if="!authenticated" type="warning" :closable="false">
		<p>You must be logged in to create content.</p>
	</el-alert>

	<space-loader v-else-if="parentId" :space-id="parentId" include-parent-path>

		<template #default="{context}">

			<form-fields
				:parent-id="parentId"
				:posting="posting"
				:context="context"
				@submit="submit"
				/>

		</template>

	</space-loader>

	<el-alert v-else type="error" :closable="false">
		<p>A parent space is required to create a link.</p>
	</el-alert>

</div>
</template>

<script>
import SpaceLoader from '@/widgets/space-loader.vue';
import FormFields from './link-form.vue';

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
			return this.$route.query.parentId ? parseInt(this.$route.query.parentId) : null;
		},
		authenticated() {
			return this.$store.getters.authenticated;
		},
	},
	methods: {
		submit(payload) {
			this.posting = true;
			ajaxPost('/ajax/space/create/link', {
				parentId: this.parentId,
				...payload,
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
