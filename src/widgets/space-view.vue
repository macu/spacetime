<template>
<div class="space">

	<check-in-button :parent-id="space.id"/>

	<horizontal-controls class="nowrap horizontal-scroll">

		<create-dropdown :parent-id="space.id"/>

		<el-input
			class="search flex-1"
			v-model="search"
			placeholder="Search in this space"
			clearable
			/>

	</horizontal-controls>

	<spaces-list
		:spaces="displaySpaces"
		@load-more="loadMode()"
		/>

	<footer class="tags-area flex-row-md" :class="{'show-all': showingAllTags}">

		<div class="collapsible flex-1 flex-row-md">

			<el-input v-if="addingTag" v-model="newTag" placeholder="Tag name">
				<template slot="append">
					<el-button @click="addingTag = false">Cancel</el-button>
					<el-button @click="addTag()">Add</el-button>
				</template>
			</el-input>
			<el-button v-else @click="addingTag = true">Add tag</el-button>

			<template v-if="tagSearchResults.length">
				<el-divider>Search results</el-divider>
				<space-tag v-for="t in tagSearchResults" :tag="t"/>
			</template>
			<template v-else>
				<space-tag v-for="t in tags" :tag="t"/>
			</template>

		</div>

		<el-button v-if="!showingAllTags">Show all tags</el-button>

	</footer>

</div>
</template>

<script>
export default {
	props: {
		spaceId: { // TODO Load space by ID
			type: Number,
			required: true,
		},
		includeTags: {
			type: Boolean,
			default: true,
		},
		includeSubspaces: {
			type: Boolean,
			default: true,
		},
	},
	data() {
		return {
			showing: 'subspaces',

			subspaces: [],

			subspaceSearchResults: [],

			addingTitle: false,
			newTitle: '',
			titleSearchResults: [],

			addingTag: false,
			newTag: '',
			tagSearchResults: [],

			allTitles: [],
			allTitlesCount: 0,

			allTags: [],
			allTagsCount: 0,
		};
	},
	computed: {
		showingAllTitles() {
			return this.showing = 'titles';
		},
		showingAllTags() {
			return this.showing = 'tags';
		},
		lastUserTitle() {
			return this.space.lastUserTitle || null;
		},
		displaySpaces() {
			return this.subspaceSearchResults.length > 0
				? this.subspaceSearchResults : this.subspaces;
		},
	},
};
</script>

<style lang="scss">
.titles-area, .tags-area {
	max-height: 80vh;
	overflow-y: auto;
}
</style>
