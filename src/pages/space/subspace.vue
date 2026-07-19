<template>
<space-output
	:space="subspace"
	:context="context"
	sub-space
	:show-reorder="showReorder"
	:data-space-id="subspace.id"
	class="space-subspace">

	<template #actions-area>
		<div class="flex-row-md">

			<el-button
				v-if="context.userAllowPinSubs"
				:type="subspace.isPinned ? 'success' : 'primary'"
				@click="togglePinned()">
				<material-icon v-if="subspace.isPinned" icon="keep"/>
				<material-icon v-else icon="keep_off"/>
				{{ subspace.isPinned ? 'Unpin' : 'Pin' }}
			</el-button>

			<checkin-button :space="subspace"/>

			<bookmark-button :space="subspace"/>

			<add-tag-button
				:parent-id="subspace.id"
				@added="tag => tags.push(tag)"
			/>

			<space-tag
				v-for="t in uniqueTags"
				:key="t.id"
				:space="t"
			/>

			<el-button
				v-if="showLoadMoreTags"
				@click="loadMoreTags()"
				type="primary">
				Load more tags
			</el-button>

		</div>
	</template>

</space-output>
</template>

<script>
import SpaceOutput from '@/widgets/space-output.vue';
import CheckinButton from '@/widgets/checkin-button.vue';
import BookmarkButton from '@/widgets/bookmark-button.vue';
import SpaceTag from '@/widgets/space-tag.vue';
import AddTagButton from '@/widgets/add-tag-button.vue';

import {
	ajaxGet,
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	emits: [
		'toggle-pinned',
	],
	components: {
		SpaceOutput,
		CheckinButton,
		BookmarkButton,
		SpaceTag,
		AddTagButton,
	},
	props: {
		subspace: {
			type: Object,
			required: true,
		},
		context: {
			type: Object,
			required: true,
		},
		showReorder: {
			type: Boolean,
			default: false,
		},
		filterJson: {
			type: String,
			required: false,
		},
	},
	data() {
		return {
			tags: this.subspace.tags || [],
			showLoadMoreTags: this.subspace.tags
				? this.subspace.tags.length === this.$const.defaultTagsLimit
				: false,
			loadingTags: false,
		};
	},
	computed: {
		allowPinning() {
			return this.context && this.context.userAllowPinSubs;
		},
		uniqueTags() {
			const ids = new Set();
			return this.tags.filter(tag => {
				if (ids.has(tag.id)) {
					return false;
				}
				ids.add(tag.id);
				return true;
			});
		},
	},
	methods: {
		togglePinned() {
			this.$emit('toggle-pinned');
		},
		loadMoreTags() {
			this.loadingTags = true;
			ajaxGet('/ajax/space/tags', {
				parentId: this.subspace.id,
				offset: this.tags.length,
				limit: this.$const.defaultTagsLimit,
				filter: this.filterJson,
			}).then(response => {
				this.tags = this.tags.concat(response);
				this.showLoadMoreTags = response.length === this.$const.defaultTagsLimit;
			}).finally(() => {
				this.loadingTags = false;
			});
		},
	},
};
</script>
