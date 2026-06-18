<template>
<form-field title="Tags" class="create-space-tags-field">

	<el-input
		v-model="addTagLabel"
		@keyup.enter="addTag"
		:disabled="disabled"
		:maxlength="$store.getters.tagMaxLength"
		show-word-limit
		placeholder="Enter tag and press enter"
		size="large"
	/>

	<div v-if="showTagsArea" class="tags flex-column">

		<div v-for="(tag, index) in pinnedTags" class="tag flex-row">
			<material-icon icon="label"/>
			<strong class="flex-1" v-text="tag"/>
			<el-button @click="unpin(index)">
				<material-icon icon="push_pin"/>
				<span>Unpin</span>
			</el-button>
			<el-button @click="moveUp(index)" :disabled="disabled || index <= 0">
				<material-icon icon="arrow_upward" />
			</el-button>
			<el-button @click="moveDown(index)" :disabled="disabled || index >= pinnedTags.length - 1">
				<material-icon icon="arrow_downward" />
			</el-button>
			<el-button @click="removePinnedTag(index)" type="warning" :disabled="disabled">
				<material-icon icon="close" />
			</el-button>
		</div>

		<div v-for="(tag, index) in tags" class="tag flex-row">
			<material-icon icon="label"/>
			<strong class="flex-1" v-text="tag"/>
			<el-button v-if="allowPinning" @click="pin(index)">
				<material-icon icon="push_pin"/>
				<span>Pin</span>
			</el-button>
			<el-button @click="removeTag(index)" type="warning" :disabled="disabled">
				<material-icon icon="close" />
			</el-button>
		</div>

	</div>

</form-field>
</template>

<script>
export default {
	emits: [
		'update:tags',
		'update:pinnedTags',
	],
	props: {
		allowPinning: {
			type: Boolean,
			default: false,
		},
		disabled: {
			type: Boolean,
			default: false,
		},
		tags: {
			type: Array,
			default: () => [],
		},
		pinnedTags: {
			type: Array,
			default: () => [],
		},
	},
	data() {
		return {
			addTagLabel: '',
		};
	},
	computed: {
		showTagsArea() {
			return this.tags.length > 0 || this.pinnedTags.length > 0;
		},
	},
	methods: {
		addTag() {
			const newTag = this.addTagLabel.trim();
			if (this.pinnedTags.includes(newTag) || this.tags.includes(newTag)) {
				return;
			}
			let updatedTags = [...this.tags, newTag];
			this.$emit('update:tags', updatedTags);
			this.addTagLabel = '';
		},
		removeTag(index) {
			let updatedTags = [...this.tags];
			updatedTags.splice(index, 1);
			this.$emit('update:tags', updatedTags);
		},
		pin(index) {
			let updatedTags = [...this.tags];
			const tag = updatedTags.splice(index, 1)[0];
			this.$emit('update:tags', updatedTags);
			this.$emit('update:pinnedTags', [...this.pinnedTags, tag]);
		},
		unpin(index) {
			let updatedPinnedTags = [...this.pinnedTags];
			const tag = updatedPinnedTags.splice(index, 1)[0];
			this.$emit('update:tags', [...this.tags, tag]);
			this.$emit('update:pinnedTags', updatedPinnedTags);
		},
		moveUp(index) {
			if (index <= 0) {
				return;
			}
			let updatedPinnedTags = [...this.pinnedTags];
			const tag = updatedPinnedTags[index];
			updatedPinnedTags.splice(index, 1);
			updatedPinnedTags.splice(index - 1, 0, tag);
			this.$emit('update:pinnedTags', updatedPinnedTags);
		},
		moveDown(index) {
			if (index >= this.pinnedTags.length - 1) {
				return;
			}
			let updatedPinnedTags = [...this.pinnedTags];
			const tag = updatedPinnedTags[index];
			updatedPinnedTags.splice(index, 1);
			updatedPinnedTags.splice(index + 1, 0, tag);
			this.$emit('update:pinnedTags', updatedPinnedTags);
		},
		removePinnedTag(index) {
			let updatedPinnedTags = [...this.pinnedTags];
			updatedPinnedTags.splice(index, 1);
			this.$emit('update:pinnedTags', updatedPinnedTags);
		},
	},
};
</script>

<style lang="scss">
.create-space-tags-field {
	.tags {
		>.tag {
			background-color: #eee;
			border-radius: 4px;
			padding: 4px 8px;
		}
	}
}
</style>
