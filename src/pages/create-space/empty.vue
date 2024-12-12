<template>
<form-layout title="Create empty space">

	<form-field title="First title" required>
		<input v-model="title" :maxlength="$store.getters.titleMaxLength" show-word-count/>
	</form-field>

	<form-actions>
		<el-button @click="create()" type="primary" :disabled="createDisabled">
			Create
		</el-button>
	</form-actions>

</form-layout>
</template>

<script>
import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	data() {
		return {
			title: '',
			posting: false,
		};
	},
	computed: {
		parentId() {
			return this.$route.params.parentId;
		},
		createDisabled() {
			return this.posting || !this.title.trim();
		},
	},
	methods: {
		create() {
			if (this.createDisabled) {
				return;
			}
			this.posting = true;
			ajaxPost('/ajax/space/create/empty', {
				parentId: this.parentId,
				title: this.title.trim(),
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
