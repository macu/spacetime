<template>
<div class="category-create-page flex-column-lg page-width-md">

	<loading-message v-if="loadingPath"/>

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
					v-model="description"
					type="textarea"
					placeholder="Description"
					:maxlength="descriptionMaxLength"
					:autosize="{minRows: 2}"
					show-word-limit
					clearable
					/>
			</div>

			<loading-message
				v-if="searchingExisting"
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

				<node-list :nodes="existing">
					<template #node-actions="{node}">
						<el-button @click="gotoNode(node)" type="primary">
							Go to category
						</el-button>
					</template>
				</node-list>

				<el-button
					type="primary" :disabled="submitDisabled" @click="create()">
					Create new category
				</el-button>

			</template>

			<el-button
				v-else
				type="primary" :disabled="submitDisabled" @click="create()">
				Create new category
			</el-button>

		</div>

	</template>

	<el-alert
		v-else-if="maxDepthExceeded"
		title="Maximum depth reached"
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
	BODY_TYPE,
} from '@/const.js';

export default {
	components: {
		ParentPath,
		NodeList,
	},
	data() {
		return {
			loadingPath: true,
			path: [],
			title: '',
			description: '',
			bodyType: BODY_TYPE.PLAINTEXT,
			loadingMarkdownPreview: false,
			markdownPreview: null,
			searchingExisting: false,
			existing: null, // null if not yet loaded
			creating: false,
		};
	},
	computed: {
		BODY_TYPE() {
			return BODY_TYPE;
		},
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
		allowMarkdownPreview() {
			return this.bodyType === BODY_TYPE.MARKDOWN && !!this.description.trim();
		},
		createAllowed() {
			return !this.maxDepthExceeded && allowCreateCategory(this.path);
		},
		submitDisabled() {
			return !this.createAllowed || !this.title.trim();
		},
		existingFound() {
			return this.existing !== null && this.existing.length > 0;
		},
	},
	watch: {
		title() {
			this.existing = null;
		},
		description() {
			this.existing = null;
			this.markdownPreview = null;
		},
	},
	beforeRouteEnter(to, from, next) {
		next(vm => {
			vm.init(to.query);
		});
	},
	beforeRouteUpdate(to, from, next) {
		this.loadingPath = true;
		this.path = [];
		this.title = '';
		this.description = '';
		this.creating = false;
		this.existing = null;
		next();
		this.init(to.query);
	},
	methods: {
		init(queryParams) {
			if (queryParams && queryParams.parentId) {
				this.loadParentPath(queryParams.parentId);
			} else {
				this.loadingPath = false;
			}
		},
		loadParentPath(parentId) {
			if (!parentId) {
				this.loadingPath = false;
				return;
			}
			this.loadingPath = true;
			ajaxGet('/ajax/node/load-path', {
				id: parentId,
			}).then(data => {
				this.path = data.path;
			}).finally(() => {
				this.loadingPath = false;
			});
		},
		create() {
			if (this.submitDisabled) {
				return;
			}
			if (this.existing === null) {
				// Search for existing
				this.searchingExisting = true;
				ajaxGet('/ajax/node/find-existing', {
					parentId: this.parentId,
					title: this.title.trim(),
					description: this.description.trim(),
					class: NODE_CLASS.CATEGORY,
				}).then(response => {
					this.existing = response.nodes;
					if (!response.nodes.length) {
						this.create();
					}
				}).finally(() => {
					this.searchingExisting = false;
				});
			} else {
				this.creating = true;
				ajaxPost('/ajax/node/create', {
					parentId: this.parentId,
					title: this.title.trim(),
					description: this.description.trim(),
					class: NODE_CLASS.CATEGORY,
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
