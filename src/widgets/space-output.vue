<template>
<div class="space-output" @click.stop="gotoSpace()">

	<div v-if="showPath && hasParentPath" @click.stop class="parent-path">
		<div
			v-for="p in space.parentPath"
			:key="p.id"
			@click.stop="gotoSpace(p)"
			class="flex-row-md">

			<material-icon icon="arrow_right_alt"/>

			<space-type :type="p.spaceType"/>

			<space-title
				v-if="p.spaceType === SPACE_TYPES.TITLE"
				:space="p"
				:show-checkin="false"
				/>

			<space-tag
				v-else-if="p.spaceType === SPACE_TYPES.TAG"
				:space="p"
				/>

			<template v-else-if="p.topTitles && p.topTitles.length > 0">
				<space-title
					v-for="title in p.topTitles"
					:space="title"
					:show-checkin="false"
					/>
			</template>

			<space-creator
				:space="p"
				/>

		</div>
	</div>

	<div class="container flex-column-md">

		<div class="space-info-bar flex-row-md" @click.stop>
			<space-type :type="space.spaceType" @click="gotoSpace()"/>
			<el-button @click="toggleBookmark()"
				:type="isBookmarked ? 'success' : 'default'" size="small">
				<material-icon v-if="isBookmarked" icon="bookmark_border"/>
				<material-icon v-else icon="bookmark"/>
			</el-button>
			<checkin-button :space="space" size="small"/>
			<space-creator :space="space"/>
			<div class="align-end flex-row-md">
				<el-button v-if="!showTitles" @click="expandTitles = true" size="small">
					Show titles
				</el-button>
				<el-button v-if="!showTags" @click="expandTags = true" class="align-end" size="small">
					Show tags
				</el-button>
			</div>
		</div>

		<div v-if="showTitles" class="space-titles-bar flex-row-md" @click.stop>
			<strong class="label">Title(s)</strong>
			<add-title
				:parent-id="space.id"
				@added="titleSpace => userTitleAdded(titleSpace)"
				@update:adding="addingTitle = $event"
				:class="{'flex-100': addingTitle}"
				/>
			<space-title
				v-if="space.userTitle"
				:space="space.userTitle"
				ellipsis
				label="(Your title)"
				/>
			<space-title
				v-if="space.topTitle"
				:space="space.topTitle"
				ellipsis
				label="(Top title)"
				/>
			<space-title
				v-if="space.originalTitle"
				:space="space.originalTitle"
				ellipsis
				label="(Original title)"
				/>
			<space-title
				v-for="title in newTitles"
				:space="title"
				@click-title="gotoSpace(title)"
				ellipsis
				label="(New title)"
				/>
			<template v-if="showAllTitles">
				<space-title
					v-for="title in extraTitles"
					:space="title"
					@click-title="gotoSpace(title)"
					ellipsis
					/>
			</template>
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
				ellipsis
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

import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	name: 'space-output', // recursive
	components: {
		CheckinButton,
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
			isBookmarked: this.space.userBookmark || false,

			addingTitle: false,
			newTitles: [],
			expandTitles: false,
			showAllTitles: false,
			extraTitles: [],

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
		showTitles() {
			if (this.expandTitles) {
				return true;
			}
			switch (this.space.spaceType) {
				case SPACE_TYPES.CHECK_IN:
				case SPACE_TYPES.TITLE:
				case SPACE_TYPES.TAG:
					return false;
			}
			// All other types show by default
			return true;
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
			switch (this.space.spaceType) {
				case SPACE_TYPES.CHECK_IN:
				case SPACE_TYPES.TITLE:
				case SPACE_TYPES.TAG:
					return false;
			}
			// All other types show by default
			return true;
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
		toggleBookmark() {
			let newState = !this.isBookmarked;
			this.isBookmarked = newState;
			ajaxPost('/ajax/bookmark', {
				spaceId: this.space.id,
				bookmark: newState,
			});
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

		>.space-info-bar, >.space-titles-bar, >.space-tags-bar {
			cursor: default;
			font-size: smaller;
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
		>.space-title, >.space-tag {
			padding: 20px;
		}
		>.space-text {
			padding: 40px;
		}
		>.portal {
			padding: 20px;
		}
	}
}
</style>
