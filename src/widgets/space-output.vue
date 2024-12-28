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

		<div @click.stop class="space-info-bar flex-row-md">
			<space-type :type="space.spaceType"/>
			<checkin-button :space="space"/>
			<space-creator :space="space"/>
			<div class="align-end flex-row-md">
				<el-button v-if="!showTitles" @click="expandTitles = true">
					Show titles
				</el-button>
				<el-button v-if="!showTags" @click="expandTags = true" class="align-end">
					Show tags
				</el-button>
			</div>
		</div>

		<div v-if="showTitles" @click.stop class="space-title-bar flex-column-sm">
			<div class="label">Title(s)</div>
			<div :class="addingTitle ? 'flex-column' : ['flex-row', 'nowrap']">
				<add-title
					:parent-id="space.id"
					@added="titleSpace => userTitleAdded(titleSpace)"
					@update:adding="addingTitle = $event"
					/>
				<div class="horizontal-scroll">
					<div class="flex-row-md nowrap">
						<space-title
							v-for="title in titlesToShow"
							:space="title"
							@click-title="gotoSpace(title)"
							/>
						<el-button>View all</el-button>
					</div>
				</div>
			</div>
		</div>

		<div v-if="showTags" @click.stop class="space-tags-bar flex-column-sm">
			<div class="label">Tag(s)</div>
			<div :class="addingTag ? 'flex-column' : ['flex-row', 'nowrap']">
				<add-tag
					:parent-id="space.id"
					@added="tagSpace => userTagAdded(tagSpace)"
					@update:adding="addingTag = $event"
					/>
				<div class="horizontal-scroll">
					<div class="flex-row-md nowrap">
						<space-tag
							v-for="tag in tagsToShow"
							:space="tag"
							@click-tag="gotoSpace(tag)"
							/>
						<el-button>View all</el-button>
					</div>
				</div>
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
	},
};
</script>

<style lang="scss">
$border-radius: 12px;

.space-output {
	border-radius: $border-radius;
	box-shadow: 0 0 5px 0 rgba(255, 255, 255);

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
			background-color: rgb(200, 200, 255);
			border-radius: $border-radius;
			padding: 10px 20px;
			cursor: default;
		}

		>.space-title-bar {
			background-color: rgb(100, 100, 200);
			border-radius: $border-radius;
			padding: 10px 20px;
			cursor: default;

			.label {
				color: white;
			}
		}

		>.space-tags-bar {
			background-color: rgb(105, 19, 98);
			border-radius: $border-radius;
			padding: 10px 20px;
			cursor: default;

			.label {
				color: white;
			}
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
			padding: 20px;
			background-color: #03d1ff;
			border-radius: 12px;
			cursor: default;
		}

		>.space-info-bar, >.space-title-bar, >.space-tags-bar, >.portal {
			// inner drop shadow
			box-shadow: inset 0 0 10px 0 rgba(255, 255, 255);
		}
	}

	>.parent-path + .container {
		border-top-left-radius: 0;
		border-top-right-radius: 0;
	}

}
</style>
