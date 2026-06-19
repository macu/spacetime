<template>
<div class="create-text-page flex-column-lg page-width-md">

	<space-loader v-if="parentId" :space-id="parentId" include-parent-path>

		<template #default="{permissions}">

			<form-fields
				:posting="posting"
				:initial-save-recording="initialSaveRecording"
				:permissions="permissions"
				@submit="submit"
				/>

		</template>

	</space-loader>

	<el-alert v-else type="error" :closable="false">
		<p>A parent space is required to create a text.</p>
	</el-alert>

</div>
</template>

<script>
import SpaceLoader from '@/widgets/space-loader.vue';
import FormFields from './text-form.vue';

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
			return this.$route.query.parentId || null;
		},
		queryType() {
			return this.$route.query.t || null;
		},
		initialSaveRecording() {
			return this.queryType === 'n';
		},
	},
	methods: {
		submit(payload) {
			this.posting = true;
			ajaxPost('/ajax/space/create/text', {
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
