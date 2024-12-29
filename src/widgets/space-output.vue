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

		<div @click.stop class="space-info-bar flex-column-md">

			<div class="flex-row-md">
				<space-type :type="space.spaceType"/>
				<checkin-button :space="space"/>
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

			<div v-if="showTitles" class="flex-row">
				<strong class="label">Title(s)</strong>
				<add-title
					:parent-id="space.id"
					@added="titleSpace => userTitleAdded(titleSpace)"
					@update:adding="addingTitle = $event"
					:class="{'flex-100': addingTitle}"
					/>
				<space-title
					v-for="title in titlesToShow"
					:space="title"
					@click-title="gotoSpace(title)"
					/>
				<el-button size="small">Load more</el-button>
			</div>

			<div v-if="showTags" class="flex-row">
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
			addingTitle: false,
			userTitles: this.space.userTitle ? [this.space.userTitle] : [],
			expandTitles: false,

			addingTag: false,
			userTags: [],
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
		titlesToShow() {
			let all = this.userTitles.concat(this.topTitles);
			return all.filter((t, i) => all.findIndex(t2 => t2.id === t.id) === i);
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
			let all = this.userTags.concat(this.topTags);
			return all.filter((t, i) => all.findIndex(t2 => t2.id === t.id) === i);
		},
	},
	methods: {
		userTitleAdded(titleSpace) {
			this.userTitles.unshift(titleSpace); // add to start
		},
		userTagAdded(tagSpace) {
			this.userTags.unshift(tagSpace); // add to start
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
	}

	>.container {
		border: medium solid darkblue;
		background-color: black;
		border-radius: $border-radius;
		padding: 20px;
		cursor: pointer; // clickable spaces

		>.space-info-bar {
			border-radius: $border-radius;
			padding: 10px 20px;
			cursor: default;

			background-color: $space-info-bar-bg-color;
			color: $space-info-bar-color;

			.space-title .text, .space-tag .text {
				white-space: nowrap;
			}
		}

		>.space-title-bar {
			padding: 10px 20px;
			cursor: default;

			background-color: $space-titles-bg-color;
			color: $space-titles-color;

			.nowrap .space-title .text {
				white-space: nowrap;
			}
		}

		>.space-tags-bar {
			padding: 10px 20px;
			cursor: default;

			background-color: $space-tags-bg-color;
			color: white;

			.nowrap .space-tag .text {
				white-space: nowrap;
			}
		}

		>.space-title, >.space-tag {
			padding: 40px;
			cursor: default;
			background-color: white;
		}
		>.space-text {
			padding: 80px;
			cursor: default;
			background-color: white;
		}

		>.portal {
			padding: 40px 20px;
			background-color: $space-bg-color;
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
</style>
