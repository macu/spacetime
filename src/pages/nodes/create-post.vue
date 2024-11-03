<template>
<div class="post-edit-page flex-column-lg page-width-md">

	<loading-message v-if="initializing"/>

	<template v-else-if="createAllowed">

		<div v-if="path.length" class="flex-column">
			<strong>Path</strong>
			<parent-path :path="path"/>
		</div>

		<form-layout title="Create post">

			<form-field title="Title" required>
				<el-input
					v-model="title"
					placeholder="Enter title"
					:maxlength="titleMaxLength"
					show-word-limit
					clearable
					/>
			</form-field>

			<form-field title="Body" required>
				<div v-for="b in blocks" class="edit-post-block flex-row">
					<el-input
						v-model="b.text"
						type="textarea"
						:autosize="{minRows: 3}"
						placeholder="Enter text"
						:maxlength="blockMaxLength"
						show-word-limit
						resize="none"
						class="flex-1"
						/>
					<div v-if="blocks.length > 1" class="flex-column flex-align-end">
						<div class="flex-row">
							<el-button @click="moveBlockUp(b)" circle>
								<material-icon icon="arrow_upward"/>
							</el-button>
							<el-button @click="moveBlockDown(b)" circle>
								<material-icon icon="arrow_downward"/>
							</el-button>
						</div>
						<el-button @click="confirmDeleteBlock(b)" type="warning" circle>
							<material-icon icon="delete"/>
						</el-button>
					</div>
				</div>
				<div class="center">
					<el-button @click="addBlock()" type="primary" :disabled="addBlockDisabled">
						<material-icon icon="add"/>
						<span>Add block</span>
					</el-button>
				</div>
			</form-field>

			<form-field title="Language" required>
				<lang-select v-model="langNodeId"/>
			</form-field>

			<loading-message
				v-if="creating"
				message="Creating post..."
				/>

			<form-actions v-else>
				<el-button
					type="primary" :disabled="submitDisabled" @click="create()">
					Create post
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
		title="Creating posts is not allowed here."
		type="error"
		:closable="false"
		/>

</div>
</template>

<script>
import ParentPath from '@/widgets/parent-path.vue';
import LangSelect from '@/widgets/lang-select.vue';

import {
	ajaxGet,
	ajaxPost,
} from '@/utils/ajax.js';

import {
	allowCreatePost,
} from '@/utils/tree.js';

import {
	NODE_CLASS,
} from '@/const.js';

export default {
	components: {
		ParentPath,
		LangSelect,
	},
	data() {
		return {
			initializing: true,
			path: [],
			createAllowed: false,

			title: '',
			blocks: [],
			langNodeId: null,

			creating: false,
		};
	},
	computed: {
		parentId() {
			return this.$route.query.parentId || null;
		},
		titleMaxLength() {
			return this.$store.getters.maxLengths.postTitle;
		},
		blockMaxLength() {
			return this.$store.getters.maxLengths.postBlock;
		},
		blockMaxCount() {
			return this.$store.getters.maxLengths.postBlockCount;
		},
		addBlockDisabled() {
			return this.blocks.length >= this.blockMaxCount;
		},
		maxDepthExceeded() {
			return this.path.length >= this.$store.getters.treeMaxDepth;
		},
		hasTitle() {
			return !!this.title.trim();
		},
		hasBodyContent() {
			return this.blocks.some(block => !!block.text.trim());
		},
		submitDisabled() {
			return !this.createAllowed || this.creating ||
				!this.hasTitle || !this.hasBodyContent || !this.langNodeId;
		},
	},
	beforeRouteEnter(to, from, next) {
		next(vm => {
			vm.init(to.query);
		});
	},
	beforeRouteUpdate(to, from, next) {
		this.initializing = true;
		this.path = [];
		this.createAllowed = false;
		this.title = '';
		this.blocks = [];
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
				class: NODE_CLASS.POST,
			}).then(data => {
				this.path = data.path;
				this.createAllowed = data.createAllowed;
				if (this.createAllowed) {
					// Add default text block
					this.blocks = [{type: 'text', text: ''}];
				}
			}).finally(() => {
				this.initializing = false;
			});
		},

		addBlock() {
			if (this.addBlockDisabled) {
				return;
			}
			this.blocks.push({type: 'text', text: ''});
		},
		moveBlockUp(block) {
			const index = this.blocks.indexOf(block);
			if (index > 0) {
				this.blocks.splice(index, 1);
				this.blocks.splice(index - 1, 0, block);
			}
		},
		moveBlockDown(block) {
			const index = this.blocks.indexOf(block);
			if (index < this.blocks.length - 1) {
				this.blocks.splice(index, 1);
				this.blocks.splice(index + 1, 0, block);
			}
		},
		confirmDeleteBlock(block) {
			this.$confirm('Delete this block?', 'Confirm', {
				confirmButtonText: 'Delete',
				cancelButtonText: 'Cancel',
				type: 'warning',
			}).then(() => {
				this.blocks = this.blocks.filter(b => b !== block);
			}).catch(() => {});
		},

		create() {
			if (this.submitDisabled) {
				return;
			}
			this.creating = true;
			ajaxPost('/ajax/node/create', {
				parentId: this.parentId,
				class: NODE_CLASS.POST,
				langNodeId: this.langNodeId,
				content: JSON.stringify({
					title: this.title.trim(),
					blocks: this.blocks,
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
			if (!confirmed && (this.hasTitle || this.hasBodyContent)) {
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

<style lang="scss">
.post-edit-page {
	.edit-post-block {
		padding: 10px;
		background-color: white;
	}
}
</style>
