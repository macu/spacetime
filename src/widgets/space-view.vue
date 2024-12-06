<template>
<div class="space">

	<check-in-button :parent-id="space.id"/>

	<header class="titles-area flex-row-md" :class="{'show-all': showingAllTitles}">

		<div class="flex-1 flex-row-md nowrap horizontal-scroll">

			<el-input v-if="addingTitle" v-model="newTitle" placeholder="Title">
				<template slot="append">
					<el-button @click="addingTitle = false">Cancel</el-button>
					<el-button @click="addTitle()">Add</el-button>
				</template>
			</el-input>
			<el-button v-else @click="addingTitle = true">Add title</el-button>

			<template v-if="titleSearchResults.length">
				<el-divider>Search results</el-divider>
				<space-title v-for="t in titleSearchResults" :title="t"/>
			</template>
			<template v-else>
				<space-title v-if="lastUserTitle" :title="lastUserTitle"/>
				<space-title v-for="t in space.top_titles" :title="t"/>
			</template>

		</div>

		<el-button v-if="!showingAllTitles">Show all titles</el-button>

	</header>

	<horizontal-controls class="nowrap horizontal-scroll">

		<el-button>Create empty space</el-button>
		<el-button>Create </el-button>

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
		space: Object,
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
