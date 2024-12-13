<template>
<form-layout title="Create tag">

	<template v-if="parentId">

		<parent-titles :parent-id="parentId"/>

		<hr/>

		<form-field title="New tag" required>
			<input v-model="tag" :maxlength="$store.getters.tagMaxLength" show-word-count/>
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
import ParentTitles from '@/widgets/parent-space-titles.vue';

import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	data() {
		return {
			tag: '',
			posting: false,
		};
	},
	computed: {
		parentId() {
			return this.$route.params.parentId;
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
