<template>
<div class="comment-edit-page flex-column-lg page-width-md">

	<loading-message v-if="initializing"/>

	<template v-else-if="createAllowed">

		<node-header v-if="parent" :node="parent" show-all/>

		<form-layout title="Create comment">

			<form-field>
				<el-input
					v-model="text"
					type="textarea"
					:autosize="{minRows: 3}"
					placeholder="Enter comment"
					:maxlength="commentMaxLength"
					show-word-limit
					resize="none"
					clearable
					/>
			</form-field>

			<form-field title="Language" required>
				<lang-select v-model="langNodeId"/>
			</form-field>

			<loading-message
				v-if="creating"
				message="Creating comment..."
				/>

			<form-actions v-else>
				<el-button
					type="primary" :disabled="submitDisabled" @click="create()">
					Create comment
				</el-button>

				<el-button
					@click="cancel()"
					type="warning">
					Cancel
				</el-button>
			</form-actions>

		</form-layout>

	</template>

	<el-alert
		v-else-if="maxDepthExceeded"
		title="Maximum tree depth reached."
		type="error"
		:closable="false"
		/>

	<el-alert v-else
		title="Creating comments is not allowed here."
		type="error"
		:closable="false"
		/>

</div>
</template>

<script>
import NodeHeader from '@/widgets/node-header.vue';
import LangSelect from '@/widgets/lang-select.vue';

import {
	ajaxGet,
	ajaxPost,
} from '@/utils/ajax.js';

import {
	NODE_CLASS,
	OWNER_TYPE,
} from '@/const.js';

export default {
	components: {
		NodeHeader,
		LangSelect,
	},
	data() {
		return {
			initializing: true,
			parent: null,
			createAllowed: false,

			text: '',
			langNodeId: null,

			creating: false,
		};
	},
	computed: {
		parentId() {
			return this.$route.query.parentId || null;
		},
		commentMaxLength() {
			return this.$store.getters.maxLengths.commentBody;
		},
		maxDepthExceeded() {
			return this.parent
				? this.parent.path.length >= this.$store.getters.treeMaxDepth
				: false;
		},
		hasText() {
			return !!this.text.trim();
		},
		submitDisabled() {
			return !this.createAllowed || this.creating ||
				!this.hasText || !this.langNodeId;
		},
	},
	beforeRouteEnter(to, from, next) {
		next(vm => {
			vm.init(to.query);
		});
	},
	beforeRouteUpdate(to, from, next) {
		this.initializing = true;
		this.parent = null;
		this.createAllowed = false;
		this.text = '';
		this.creating = false;
		next();
		this.init(to.query);
	},
	methods: {
		init(queryParams) {
			if (queryParams && queryParams.parentId) {
				this.loadCreate(queryParams.parentId);
			} else {
				this.initializing = false;
			}
		},
		loadCreate(parentId) {
			if (!parentId) {
				this.initializing = false;
				this.createAllowed = true;
				return;
			}
			this.initializing = true;
			ajaxGet('/ajax/node/load-create', {
				parentId,
				class: NODE_CLASS.COMMENT,
			}).then(data => {
				this.parent = data.parent;
				this.createAllowed = data.createAllowed;
			}).finally(() => {
				this.initializing = false;
			});
		},

		create() {
			if (this.submitDisabled) {
				return;
			}
			this.creating = true;
			ajaxPost('/ajax/node/create', {
				parentId: this.parentId,
				ownerType: OWNER_TYPE.USER,
				class: NODE_CLASS.COMMENT,
				langNodeId: this.langNodeId,
				content: JSON.stringify({
					text: this.text.trim(),
				}),
			}).then(response => {
				this.$router.replace({
					name: 'node-view',
					params: {
						id: response.id,
					},
				});
			}).finally(() => {
				this.creating = false;
			});
		},

		cancel(confirmed = false) {
			if (!confirmed && this.hasText) {
				this.$confirm('Are you sure you want to cancel?', 'Unsaved changes', {
					confirmButtonText: 'Yes',
					cancelButtonText: 'No',
					type: 'warning',
				}).then(() => {
					this.cancel(true);
				});
				return;
			}

			if (this.parentId) {
				this.$router.push({
					name: 'node-view',
					params: {
						id: this.parentId,
					},
				});
			} else {
				this.$router.push({
					name: 'dashboard',
				});
			}
		},
	},
};
</script>
