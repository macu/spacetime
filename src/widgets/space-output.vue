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
			<el-button v-if="!showTitles" @click="expandTitles = true" class="align-end">
				Show titles
			</el-button>
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

		<slot name="subspaces"/>

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
			userTitles: (this.space.userTitles || []).slice(),
			expandTitles: false,
		};
	},
	computed: {
		SPACE_TYPES() {
			return SPACE_TYPES;
		},
		hasUserTitles() {
			return this.userTitles.length > 0;
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
	},
	methods: {
		userTitleAdded(titleSpace) {
			let index = this.userTitles.findIndex(t => t.id === titleSpace.id);
			if (index >= 0) {
				// Title already present
				this.userTitles.splice(index, 1); // remove
			}
			this.userTitles.unshift(titleSpace); // add to start
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

		>.space-title, >.space-tag, >.space-text {
			padding: 40px;
			cursor: default;
		}
	}

	>.parent-path + .container {
		border-top-left-radius: 0;
		border-top-right-radius: 0;
	}

}
</style>
