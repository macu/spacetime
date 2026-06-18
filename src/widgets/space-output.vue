<template>
<div class="space-output" @click.stop="gotoSpace()">

	<div v-if="showPath && hasParentPath" @click.stop class="parent-path">
		<div
			v-for="p in space.parentPath"
			:key="p.id"
			@click.stop="gotoSpace(p)"
			class="flex-row-md">

			<material-icon icon="arrow_right_alt"/>

			<space-type :space="p" show-pin/>

			<span v-if="p.spaceType === SPACE_TYPES.BRANCH"
				v-text="p.label"
				class="label"
				/>

			<span
				v-else-if="p.spaceType === SPACE_TYPES.TAG"
				v-text="p.text"
				class="label"
				/>

			<space-creator
				:space="p"
				/>

		</div>
	</div>

	<div class="container flex-column-md">

		<div class="space-info-bar flex-row-md" @click.stop>
			<el-button v-if="userAllowPin" @click="togglePinned()">
				<material-icon v-if="isPinned" icon="keep"/>
				<material-icon v-else icon="keep_off"/>
				{{ isPinned ? 'Unpin' : 'Pin' }}
			</el-button>
			<space-type :space="space" :show-pin="!userAllowPin" @click="gotoSpace()"/>
			<bookmark-button :space="space"/>
			<checkin-button :space="space"/>
			<div v-if="space.label" class="space-label">
				<strong v-text="space.label"/>
			</div>
			<space-creator :space="space"/>
			<div class="align-end flex-row-md">
				<el-button v-if="!showTags" @click="expandTags = true" class="align-end" size="small">
					Show tags
				</el-button>
			</div>
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

		<space-tag
			v-else-if="space.spaceType === SPACE_TYPES.TAG"
			:space="space"
			:show-checkin="false"
			show-pin
			/>

		<space-text
			v-else-if="space.spaceType === SPACE_TYPES.TEXT"
			:space="space"
			/>

		<div v-if="$slots.default" class="portal" @click.stop>
			<slot
				:user-allow-pin="userAllowPin"
				:user-allow-pin-subspaces-on-create="userAllowPinSubspacesOnCreate"
			/>
		</div>

	</div>

</div>
</template>

<script>
import CheckinButton from './checkin-button.vue';
import BookmarkButton from './bookmark-button.vue';
import SpaceType from './space-type.vue';
import SpaceCreator from './space-creator.vue';
import SpaceTag from './space-tag.vue';
import SpaceText from './space-text.vue';
import AddTag from './add-tag-button.vue';

import {
	SPACE_TYPES,
} from '@/const.js';

export default {
	name: 'space-output', // recursive
	emits: [
		'toggle-pinned',
	],
	components: {
		CheckinButton,
		BookmarkButton,
		SpaceType,
		SpaceCreator,
		SpaceTag,
		SpaceText,
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
		userAllowPin: {
			type: Boolean,
			default: false,
		},
		userAllowPinSubspacesOnCreate: {
			type: Function,
			required: false,
		},
	},
	data() {
		return {
			addingTag: false,
			newTags: [],
			expandTags: false,
		};
	},
	computed: {
		SPACE_TYPES() {
			return SPACE_TYPES;
		},
		hasParentPath() {
			return !!this.space.parentPath && this.space.parentPath.length > 0;
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
		isPinned() {
			return this.space.isPinned;
		},
	},
	methods: {
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
		togglePinned() {
			this.$emit('toggle-pinned', this.space);
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
		.label {
			display: inline-block;
			padding-left: 10px;
			padding-right: 10px;
			font-weight: bold;
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
