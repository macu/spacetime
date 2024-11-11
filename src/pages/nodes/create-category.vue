<template>
<div class="category-create-page flex-column-lg page-width-md">

	<loading-message v-if="initializing"/>

	<template v-else-if="createAllowed">

		<el-drawer v-if="existingFound" v-model="showingExisting">

			<template #header>
				<div class="flex-column">
					<h3>Existing categories</h3>
				</div>
			</template>

			<div class="flex-column-lg">

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

			</div>

			<template #footer>
				<el-button @click="cancel(true)" type="warning">
					<material-icon icon="cancel"/>
					<span>Cancel create category</span>
				</el-button>
			</template>

		</el-drawer>

		<node-header v-if="parent" :node="parent" show-all/>

		<form-layout title="Create category">

			<form-field title="Ownership">
				<el-radio-group v-model="ownerType" class="flex-column">
					<el-radio :label="OWNER_TYPE.PUBLIC">
						<div class="flex-row">
							<material-icon icon="public"/>
							<span>This category will belong to the public</span>
						</div>
					</el-radio>
					<el-radio :label="OWNER_TYPE.USER">
						<div class="flex-row">
							<material-icon icon="person"/>
							<span>This category will belong to me</span>
						</div>
					</el-radio>
				</el-radio-group>
				<el-alert
					v-if="ownerType === OWNER_TYPE.PUBLIC"
					type="info" effect="dark" show-icon :closable="false">
					<p>You will not be able to delete this category after it is created.</p>
					<p>Other users will be able to change the title and description.</p>
				</el-alert>
				<el-alert
					v-else-if="ownerType === OWNER_TYPE.USER"
					type="info" effect="dark" show-icon :closable="false">
					<p>You will have control over this category and the content directly within it, and your name will be displayed along with it.</p>
				</el-alert>
			</form-field>

			<form-field title="Title" required>
				<el-input
					v-model="title"
					placeholder="Title"
					:maxlength="titleMaxLength"
					show-word-limit
					clearable
					/>
			</form-field>

			<form-field title="Description">
				<el-input
					v-model="description"
					type="textarea"
					placeholder="Description"
					:maxlength="descriptionMaxLength"
					:autosize="{minRows: 2}"
					show-word-limit
					resize="none"
					clearable
					/>
			</form-field>

			<form-field title="Language" required>
				<lang-select v-model="langNodeId"/>
			</form-field>

			<loading-message
				v-if="findingExisting"
				message="Searching for existing categories..."
				/>

			<loading-message
				v-else-if="creating"
				message="Creating category..."
				/>

			<template v-else>

				<el-alert
					v-if="existingNotFound"
					type="success"
					title="No similar categories were found."
					show-icon :closable="false"
					/>

				<form-actions>
					<el-button
						@click="create()"
						type="primary"
						:disabled="submitDisabled">
						Create category
					</el-button>

					<el-button
						v-if="existingNotChecked"
						@click="findExisting()"
						:disabled="submitDisabled">
						<material-icon icon="search"/>
						<span>Check for similar existing categories</span>
					</el-button>

					<el-button
						v-else-if="existingFound"
						@click="showExisting()">
						<material-icon icon="search"/>
						<span>Show existing categories</span>
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
		title="Creating categories is not allowed here."
		type="error"
		:closable="false"
		/>

</div>
</template>

<script>
import NodeHeader from '@/widgets/node-header.vue';
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
		NodeHeader,
		NodeList,
		LangSelect,
	},
	data() {
		return {
			initializing: true,
			parent: null,
			createAllowed: false,

			ownerType: OWNER_TYPE.PUBLIC,
			title: '',
			description: '',
			langNodeId: null,

			findingExisting: false,
			existingNodes: null, // null if not yet loaded
			showingExisting: false,

			creating: false,
		};
	},
	computed: {
		OWNER_TYPE() {
			return OWNER_TYPE;
		},
		parentId() {
			return this.$route.query.parentId || null;
		},
		titleMaxLength() {
			return this.$store.getters.maxLengths.categoryTitle;
		},
		descriptionMaxLength() {
			return this.$store.getters.maxLengths.categoryDescription;
		},
		maxDepthExceeded() {
			return this.parent
				? this.parent.path.length >= this.$store.getters.treeMaxDepth
				: false;
		},
		hasTitle() {
			return !!this.title.trim();
		},
		hasDescription() {
			return !!this.description.trim();
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
		description() {
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
		this.parent = null;
		this.createAllowed = false;
		this.ownerType = OWNER_TYPE.PUBLIC;
		this.title = '';
		this.description = '';
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
				class: NODE_CLASS.CATEGORY,
			}).then(data => {
				this.parent = data.parent;
				this.createAllowed = data.createAllowed;
				if (data.parent && data.parent.path.length > 0) {
					// Default to user-owned if creating within user-owned parent category
					let parentNode = data.parent.path[data.parent.path.length - 1];
					if (
						parentNode.ownerType === OWNER_TYPE.USER &&
						parentNode.creator.id === this.$store.getters.currentUserId
					) {
						this.ownerType = OWNER_TYPE.USER;
					}
				}
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
				class: NODE_CLASS.CATEGORY,
				query: (this.title.trim() + ' ' + this.description.trim()).trim(),
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
				ownerType: this.ownerType,
				class: NODE_CLASS.CATEGORY,
				langNodeId: this.langNodeId,
				content: JSON.stringify({
					title: this.title.trim(),
					description: this.description.trim(),
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
			if (!confirmed && (this.hasTitle || this.hasDescription)) {
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
