<template>
<form-layout title="Create text" class="create-text-page page-width-md">

	<template v-if="parentId">
		<parent-titles :parent-id="parentId"/>
		<hr/>
	</template>

	<form-field title="Title">
		<el-input
			v-model="title"
			:maxlength="$store.getters.titleMaxLength"
			show-word-count
			size="large"
			/>
	</form-field>

	<form-field title="Text" required>
		<el-input
			type="textarea"
			v-model="text"
			:maxlength="$store.getters.textMaxLength"
			show-word-count
			:autosize="{minRows: 3}"
			/>
	</form-field>

	<form-actions>
		<el-button @click="create()" type="primary" :disabled="createDisabled">
			Create
		</el-button>
	</form-actions>

</form-layout>
</template>

<script>
import ParentTitles from '@/widgets/parent-space-titles.vue';

import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	components: {
		ParentTitles,
	},
	data() {
		return {
			title: '',
			text: '',
			posting: false,
		};
	},
	computed: {
		parentId() {
			return this.$route.query.parentId || null;
		},
		createDisabled() {
			return this.posting || !this.text.trim();
		},
	},
	methods: {
		create() {
			if (this.createDisabled) {
				return;
			}
			this.posting = true;
			ajaxPost('/ajax/space/create/text', {
				parentId: this.parentId,
				title: this.title.trim(),
				text: this.text.trim(),
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
