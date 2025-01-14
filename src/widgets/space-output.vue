<template>
<div class="space-output" @click.stop="gotoSpace()">

	<div v-if="showPath && hasParentPath" @click.stop class="parent-path">
		<div
			v-for="p in space.parentPath"
			:key="p.id"
			@click.stop="gotoSpace(p)"
			class="flex-row-md">

			<material-icon icon="arrow_right_alt"/>

			<space-type :space="p"/>

			<space-title
				v-if="p.spaceType === SPACE_TYPES.TITLE"
				:space="p"
				:show-checkin="false"
				/>

			<space-tag
				v-else-if="p.spaceType === SPACE_TYPES.TAG"
				:space="p"
				/>

			<space-title
				v-else-if="p.originalTitle"
				:space="p.originalTitle"
				:show-checkin="false"
				/>

			<space-creator
				:space="p"
				/>

		</div>
	</div>

	<div class="container flex-column-md">

		<div class="space-info-bar flex-row-md" @click.stop>
			<space-type :space="space" @click="gotoSpace()"/>
			<bookmark-button :space="space"/>
			<checkin-button :space="space"/>
			<space-title
				v-if="!expandTitles && firstTitle"
				:space="firstTitle"
				:label="firstTitle.label"
				/>
			<space-creator :space="space"/>
			<div class="align-end flex-row-md">
				<el-button v-if="!expandTitles" @click="expandTitles = true" size="small">
					Show titles
				</el-button>
				<el-button v-if="!showTags" @click="expandTags = true" class="align-end" size="small">
					Show tags
				</el-button>
			</div>
		</div>

		<div v-if="expandTitles" class="space-titles-bar flex-row-md" @click.stop>
			<strong class="label">Title(s)</strong>
			<add-title
				:parent-id="space.id"
				@added="titleSpace => userTitleAdded(titleSpace)"
				@update:adding="addingTitle = $event"
				:class="{'flex-100': addingTitle}"
				/>
			<space-title
				v-for="title in titles"
				:space="title"
				@click-title="gotoSpace(title)"
				:label="title.label"
				/>
			<el-button size="small" @click="showAllTitles = true">Load more</el-button>
		</div>

		<div v-if="showTags" class="space-tags-bar flex-row-md" @click.stop>
			<strong class="label">Tag(s)</strong>
			<add-tag
				:parent-id="space.id"
				@added="tagSpace => userTagAdded(tagSpace)"
				@update:adding="addingTag = $event"
				:class="{'flex-100': addingTag}"
				/>
			<space-tag
				v-for="tag in tagsToShow"
				:space="tag"
				@click-tag="gotoSpace(tag)"
				/>
			<el-button size="small">Load more</el-button>
		</div>

		<space-output
			v-if="space.spaceType === SPACE_TYPES.CHECK_IN && !!space.checkinSpace"
			:space="space.checkinSpace"
			show-path
			/>

		<space-title
			v-else-if="space.spaceType === SPACE_TYPES.TITLE"
			:space="space"
			:show-checkin="false"
			/>

		<space-tag
			v-else-if="space.spaceType === SPACE_TYPES.TAG"
			:space="space"
			:show-checkin="false"
			/>

		<space-text
			v-else-if="space.spaceType === SPACE_TYPES.TEXT"
			:space="space"
			/>

		<div v-if="$slots.default" class="portal" @click.stop>
			<slot/>
		</div>

	</div>

</div>
</template>

<script>
import CheckinButton from './checkin-button.vue';
import BookmarkButton from './bookmark-button.vue';
import SpaceType from './space-type.vue';
import SpaceCreator from './space-creator.vue';
import SpaceTitle from './space-title.vue';
import SpaceTag from './space-tag.vue';
import SpaceText from './space-text.vue';
import AddTitle from './add-title-button.vue';
import AddTag from './add-tag-button.vue';

import {
	SPACE_TYPES,
} from '@/const.js';

export default {
	name: 'space-output', // recursive
	components: {
		CheckinButton,
		BookmarkButton,
		SpaceType,
		SpaceCreator,
		SpaceTitle,
		SpaceTag,
		SpaceText,
		AddTitle,
		AddTag,
	},
	props: {
		space: {
			type: Object,
			required: true,
		},
		showPath: {
			type: Boolean,
			default: false,
		},
	},
	data() {
		return {
			addingTitle: false,
			newTitles: [],
			expandTitles: false,
			loadingTitles: false,
			moreTitles: [],

			addingTag: false,
			newTags: [],
			expandTags: false,

			titlesExpanded: false,
		};
	},
	computed: {
		SPACE_TYPES() {
			return SPACE_TYPES;
		},
		topTitles() {
			return this.space.topTitles || [];
		},
		hasTopTitles() {
			return this.topTitles.length > 0;
		},
		hasParentPath() {
			return !!this.space.parentPath && this.space.parentPath.length > 0;
		},
		titles() {
			let titles = this.newTitles.map(t => {
				return {
					...t,
					label: '(New title)',
				};
			});

			if (this.space.userTitle) {
				titles.push({
					...this.space.userTitle,
					label: '(Your title)',
				});
			}

			if (this.space.topTitle) {
				titles.push({
					...this.space.topTitle,
					label: '(Top title)',
				});
			}

			if (this.space.originalTitle) {
				titles.push({
					...this.space.originalTitle,
					label: '(Original title)',
				});
			}

			return titles.concat(this.moreTitles.map(t => {
				return {
					...t,
					label: null,
				};
			}));
		},
		firstTitle() {
			if (this.titles.length > 0) {
				return this.titles[0];
			}
			return null;
		},
		topTags() {
			return this.space.topTags || [];
		},
		hasTopTags() {
			return this.topTags.length > 0;
		},
		showTags() {
			if (this.expandTags) {
				return true;
			}
			return false;
		},
		tagsToShow() {
			let all = this.newTags.concat(this.topTags);
			return all.filter((t, i) => all.findIndex(t2 => t2.id === t.id) === i);
		},
	},
	methods: {
		userTitleAdded(titleSpace) {
			this.newTitles.unshift(titleSpace); // add to start
		},
		userTagAdded(tagSpace) {
			this.newTags.unshift(tagSpace); // add to start
		},
		gotoSpace(s = null) {
			this.$router.push({
				name: 'space',
				params: {
					spaceId: s ? s.id : this.space.id,
				},
			});
		},
		loadMoreTitles() {
			this.titlesExpanded = true;
		},
	},
};
</script>

<style lang="scss">
@import '@/styles/vars.scss';

.space-output {
	border-radius: $border-radius;
	box-shadow: $space-drop-shadow;

	>.parent-path {
		border-top-left-radius: $border-radius;
		border-top-right-radius: $border-radius;
		background-color: lightsteelblue;
		>div {
			padding: 10px 20px;
			cursor: pointer;
		}
		>div+div {
			border-top: thin solid black;
		}
		.space-title {
			padding-left: 10px;
			padding-right: 10px;
		}
	}

	>.container {
		background-color: white;
		border: thin solid darkblue;
		border-radius: $border-radius;
		padding: 20px 40px;
		cursor: pointer; // clickable spaces

		.space-type {
			cursor: pointer;
		}

		>.space-titles-bar, >.space-tags-bar {
			font-size: smaller;
			cursor: default;
		}

		>.space-title, >.space-tag {
			padding: 40px;
			cursor: default;
		}
		>.space-text {
			padding: 80px;
			cursor: default;
		}

		>.portal {
			background-color: $space-bg-color;
			border: thin solid darkblue;
			border-radius: $border-radius;
			padding: 40px;
			border-radius: 12px;
			cursor: default;

			// inner drop shadow
			box-shadow: $space-inner-drop-shadow;
		}
	}

	>.parent-path + .container {
		border-top-left-radius: 0;
		border-top-right-radius: 0;
	}

}

.is-mobile .space-output {
	>.container {
		padding: 20px 10px;
		>.space-title, >.space-tag, >.space-text {
			padding: 20px;
		}
		>.portal {
			padding: 20px;
		}
	}
}
</style>
