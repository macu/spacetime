<template>
<loading-message v-if="loading"/>
<space-output v-else-if="space" :space="space"/>
<el-alert v-else type="error" show-icon :closable="false">
	This space could not be loaded.
</el-alert>
</template>

<script>
import SpaceOutput from '@/widgets/space-output.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	data() {
		return {
			loading: true,
			space: null,
		};
	},
	computed() {
		spaceId() {
			return this.$route.params.spaceId;
		},
	},
	beforeRouteEnter(to, from, next) {
		next(vm => {
			vm.loadSpace();
		});
	},
	beforeRouteUpdate(to, from, next) {
		this.space = null;
		this.loading = true;
		next();
		this.loadSpace();
	},
	methods: {
		loadSpace() {
			this.loading = true;
			this.space = null;
			ajaxGet('/ajax/space', {
				spaceId: this.spaceId,
			}).then(response => {
				this.space = response;
			}).finally((error) => {
				this.loading = false;
			});
		},
	},
};
