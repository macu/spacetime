<template>
<div class="category-create-page flex-column-lg page-width-md">

	<loading-message v-if="initializing"/>

	<template v-else-if="createAllowed">

		<div v-if="path.length" class="flex-column">
			<strong>Path</strong>
			<parent-path :path="path"/>
		</div>

		<h2>Create category</h2>

		<div class="form-layout">

			<div class="field">
				<label>Title (required)</label>
				<el-input
					v-model="title"
					placeholder="Title"
					:maxlength="titleMaxLength"
					show-word-limit
					clearable
					/>
			</div>

			<div class="field">
				<label>Description (optional)</label>
				<el-input
					v-model="body"
					type="textarea"
					placeholder="Description"
					:maxlength="descriptionMaxLength"
					:autosize="{minRows: 2}"
					show-word-limit
					clearable
					/>
			</div>

			<loading-message
				v-if="findingExisting"
				message="Searching for existing categories..."
				/>

			<loading-message
				v-else-if="creating"
				message="Creating category..."
				/>

			<template v-else-if="existingFound">

				<el-alert
					title="Found similar existing content"
					type="success" effect="dark" show-icon :closable="false">
					Please ensure the category you are trying to create doesn't already exist.
				</el-alert>

				<node-list :nodes="existingNodes">
					<template #node-actions="{node}">
						<el-button @click="gotoNode(node)" type="primary">
							<material-icon icon="arrow_forward"/>
							<span>Go to category</span>
						</el-button>
					</template>
				</node-list>

				<div>
					<el-button
						type="primary" :disabled="submitDisabled" @click="create()">
						Create category
					</el-button>
				</div>

			</template>

			<div v-else>
				<el-button
					type="primary" :disabled="submitDisabled" @click="create()">
					Create category
				</el-button>
			</div>

		</div>

	</template>

	<el-alert
		v-else-if="maxDepthExceeded"
		title="Maximum depth reached."
		type="error"
		:closable="false"
		/>

	<el-alert v-else
		title="Creating categories is not allowed here."
		type="error"
		:closable="false"
		/>

</div>
</template>

<script>
import ParentPath from '@/widgets/parent-path.vue';
import NodeList from '@/widgets/node-list.vue';

import {
	ajaxGet,
	ajaxPost,
} from '@/utils/ajax.js';

import {
	allowCreateCategory,
} from '@/utils/tree.js';

import {
	NODE_CLASS,
} from '@/const.js';

export default {
	components: {
		ParentPath,
		NodeList,
	},
	data() {
		return {
			initializing: true,
			path: [],
			createAllowed: false,

			title: '',
			body: '',

			findingExisting: false,
			existingNodes: null, // null if not yet loaded

			creating: false,
		};
	},
	computed: {
		parentId() {
			return this.$route.query.parentId || null;
		},
		titleMaxLength() {
			return this.$store.getters.categoryTitleMaxLength;
		},
		descriptionMaxLength() {
			return this.$store.getters.categoryDescriptionMaxLength;
		},
		maxDepthExceeded() {
			return this.path.length >= this.$store.getters.treeMaxDepth;
		},
		submitDisabled() {
			return !this.createAllowed || !this.title.trim();
		},
		existingFound() {
			return this.existingNodes !== null && this.existingNodes.length > 0;
		},
	},
	watch: {
		title() {
			this.existingNodes = null;
		},
		body() {
			this.existingNodes = null;
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
		this.body = '';
		this.existingNodes = null;
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
				class: NODE_CLASS.CATEGORY,
			}).then(data => {
				this.path = data.path;
				this.createAllowed = data.createAllowed;
			}).finally(() => {
				this.initializing = false;
			});
		},
		create() {
			if (this.submitDisabled) {
				return;
			}
			if (this.existingNodes === null) {
				// Search for existing
				this.findingExisting = true;
				let query = (this.title.trim() + ' ' + this.body.trim()).trim();
				ajaxGet('/ajax/node/find-existing', {
					parentId: this.parentId,
					class: NODE_CLASS.CATEGORY,
					query,
				}).then(response => {
					this.existingNodes = response.nodes;
					if (!response.nodes.length) {
						this.create();
					}
				}).finally(() => {
					this.findingExisting = false;
				});
			} else {
				this.creating = true;
				ajaxPost('/ajax/node/create', {
					parentId: this.parentId,
					class: NODE_CLASS.CATEGORY,
					title: this.title.trim(),
					body: this.body.trim(),
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
			}
		},
		gotoNode(node) {
			this.$router.push({
				name: 'node-view',
				params: {
					id: node.id,
				},
			});
		},
	},
};
</script>
