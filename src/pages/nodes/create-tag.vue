<template>
<div class="tag-create-page flex-column-lg page-width-md">

	<loading-message v-if="initializing"/>

	<template v-else-if="createAllowed">

		<el-drawer v-if="existingFound" v-model="showingExisting" class="flex-column-drawer">

			<template #header>
				<h3>Existing tags</h3>
			</template>

			<el-alert
				title="Found similar existing content"
				type="success" effect="dark" show-icon :closable="false">
				Please ensure the tag you are trying to create doesn't already exist.
			</el-alert>

			<node-list :nodes="existingNodes">
				<template #node-actions="{node}">
					<el-button @click="gotoNode(node)" type="primary">
						<material-icon icon="arrow_forward"/>
						<span>Go to tag</span>
					</el-button>
				</template>
			</node-list>

		</el-drawer>

		<form-layout title="Create tag">

			<form-field title="Title" required>
				<el-input
					v-model="title"
					placeholder="Title"
					:maxlength="titleMaxLength"
					show-word-limit
					clearable
					/>
			</form-field>

			<form-field title="Language" required>
				<lang-select v-model="langNodeId"/>
			</form-field>

			<loading-message
				v-if="findingExisting"
				message="Searching for existing tags..."
				/>

			<loading-message
				v-else-if="creating"
				message="Creating tag..."
				/>

			<template v-else>

				<el-alert
					v-if="existingNotFound"
					type="success"
					title="No similar tags were found."
					show-icon :closable="false"
					/>

				<form-actions>
					<el-button
						@click="create()"
						type="primary"
						:disabled="submitDisabled">
						Create tag
					</el-button>

					<el-button
						v-if="existingNotChecked"
						@click="findExisting()"
						:disabled="submitDisabled">
						<material-icon icon="search"/>
						<span>Check for similar existing tags</span>
					</el-button>

					<el-button
						v-else-if="existingFound"
						@click="showExisting()">
						<material-icon icon="search"/>
						<span>Show existing tags</span>
					</el-button>

					<el-button
						@click="cancel()"
						type="warning">
						Cancel
					</el-button>
				</form-actions>

			</template>

		</form-layout>

	</template>

	<el-alert
		v-else-if="maxDepthExceeded"
		title="Maximum tree depth reached."
		type="error"
		:closable="false"
		/>

	<el-alert v-else
		title="Creating tags is not allowed here."
		type="error"
		:closable="false"
		/>

</div>
</template>

<script>
import ParentPath from '@/widgets/parent-path.vue';
import NodeList from '@/widgets/node-list.vue';
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
		ParentPath,
		NodeList,
		LangSelect,
	},
	data() {
		return {
			initializing: true,
			path: [],
			createAllowed: false,

			title: '',
			langNodeId: null,

			findingExisting: false,
			existingNodes: null, // null if not yet loaded
			showingExisting: false,

			creating: false,
		};
	},
	computed: {
		parentId() {
			return this.$route.query.parentId || null;
		},
		titleMaxLength() {
			return this.$store.getters.maxLengths.tagTitle;
		},
		maxDepthExceeded() {
			return this.path.length >= this.$store.getters.treeMaxDepth;
		},
		hasTitle() {
			return !!this.title.trim();
		},
		submitDisabled() {
			return !this.createAllowed || this.creating ||
				!this.hasTitle || !this.langNodeId;
		},
		existingNotChecked() {
			return this.existingNodes === null;
		},
		existingFound() {
			return this.existingNodes !== null && this.existingNodes.length > 0;
		},
		existingNotFound() {
			return this.existingNodes !== null && this.existingNodes.length === 0;
		},
	},
	watch: {
		title() {
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
		this.ownerType = OWNER_TYPE.PUBLIC;
		this.title = '';
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
				this.loadCreate();
			}
		},
		loadCreate(parentId = null) {
			if (!parentId) {
				this.initializing = false;
				this.createAllowed = true;
				return;
			}
			this.initializing = true;
			ajaxGet('/ajax/node/load-create', {
				parentId,
				class: NODE_CLASS.TAG,
			}).then(data => {
				this.path = data.path;
				this.createAllowed = data.createAllowed;
			}).finally(() => {
				this.initializing = false;
			});
		},
		findExisting() {
			if (this.submitDisabled) {
				return;
			}
			this.findingExisting = true;
			ajaxGet('/ajax/node/find-existing', {
				parentId: this.parentId,
				class: NODE_CLASS.TAG,
				query: this.title.trim(),
			}).then(response => {
				this.existingNodes = response.nodes;
				if (response.nodes.length > 0) {
					this.showExisting();
				}
			}).finally(() => {
				this.findingExisting = false;
			});
		},
		showExisting() {
			this.showingExisting = true;
		},
		create() {
			if (this.submitDisabled) {
				return;
			}
			this.creating = true;
			ajaxPost('/ajax/node/create', {
				parentId: this.parentId,
				ownerType: OWNER_TYPE.PUBLIC,
				class: NODE_CLASS.TAG,
				langNodeId: this.langNodeId,
				content: JSON.stringify({
					title: this.title.trim(),
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
		gotoNode(node) {
			this.$router.push({
				name: 'node-view',
				params: {
					id: node.id,
				},
			});
		},
		cancel(confirmed = false) {
			if (!confirmed && this.hasTitle) {
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
