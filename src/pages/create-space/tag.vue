<template>
<form-layout title="Create tag" class="create-tag-page page-width-md">

	<template v-if="parentId">

		<parent-path :parent-id="parentId"/>

		<hr/>

		<form-field title="Name for new tag" required>
			<el-input
				v-model="tag"
				:maxlength="$store.getters.tagMaxLength"
				show-word-limit
				size="large"
				/>
		</form-field>

		<form-actions>
			<el-button @click="create()" type="primary" :disabled="createDisabled">
				Create
			</el-button>
		</form-actions>

	</template>

	<el-alert v-else type="error" :closable="false">
		<p>A parent space is required to create a tag.</p>
	</el-alert>

</form-layout>
</template>

<script>
import ParentPath from '@/widgets/parent-path.vue';

import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	components: {
		ParentPath,
	},
	data() {
		return {
			tag: '',
			posting: false,
		};
	},
	computed: {
		parentId() {
			return this.$route.query.parentId || null;
		},
		createDisabled() {
			return this.posting || !this.tag.trim();
		},
	},
	methods: {
		create() {
			if (this.createDisabled) {
				return;
			}
			this.posting = true;
			ajaxPost('/ajax/space/create/tag', {
				parentId: this.parentId,
				tag: this.tag.trim(),
			}).then(response => {
				this.$router.replace({
					name: 'space',
					params: {
						spaceId: response.id,
					},
				});
			}).catch(() => {
				this.posting = false;
			})
		},
	},
};
</script>
