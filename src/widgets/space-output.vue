<template>
<div class="space-output flex-column" @click="gotoSpace()">

	<div @click.stop class="space-info-bar flex-row-md">
		<checkin-button :space="space"/>
		<space-type :type="space.spaceType"/>
	</div>

	<div v-if="hasUserTitles" @click.stop class="space-title-bar flex-row-md nowrap">
		<add-title @added="titleSpace => showUserTitle(titleSpace)"/>
		<div class="flex-row-md nowrap horizontal-scroll">
			<space-title v-for="title in userTitles" :space="title"/>
			<el-button>View all</el-button>
		</div>
	</div>

	<div @click.stop class="space-title-bar flex-row-md nowrap">
		<add-title v-if="!hasUserTitles" @added="titleSpace => showUserTitle(titleSpace)"/>
		<div v-if="hasTopTitles" class="flex-row-md nowrap horizontal-scroll">
			<space-title v-for="title in userTitles" :space="title"/>
			<el-button>View all</el-button>
		</div>
	</div>

<!--
	<space-output
		v-if="spaceType === SPACE_TYPES.CHECKIN && !!space.checkinSpace"
		:space="space.checkinSpace"
		/>

	<space-title
		v-else-if="spaceType === SPACE_TYPES.TITLE"
		:title-space="space"
		/>

	<space-tag
		v-else-if="spaceType === SPACE_TYPES.TAG"
		:space="space"
		/>

	<space-text
		v-else-if="spaceType === SPACE_TYPES.TEXT"
		:space="space"
		/> -->

</div>
</template>

<script>
import CheckinButton from './checkin-button.vue';
import SpaceType from './space-type.vue';
import SpaceTitle from './space-title.vue';
import SpaceText from './space-text.vue';
import AddTitle from './add-space-title.vue';

import {
	SPACE_TYPES,
} from '@/const.js';

export default {
	name: 'space-output', // recursive
	components: {
		CheckinButton,
		SpaceType,
		SpaceTitle,
		SpaceText,
		AddTitle,
	},
	props: {
		space: {
			type: Object,
			required: true,
		},
	},
	data() {
		return {
			userTitles: (this.space.userTitles || []).slice(),
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
	},
	methods: {
		showUserTitle(titleSpace) {
			let index = this.userTitles.findIndex(t => t.id === titleSpace.id);
			if (index >= 0) {
				// Title already present
				this.userTitles.splice(index, 1); // remove
			}
			this.userTitles.unshift(titleSpace); // add to start
		},
		gotoSpace() {
			this.$router.push({
				name: 'space',
				params: {
					spaceId: this.space.id,
				},
			});
		},
	},
};
</script>

<style lang="scss">
$border-radius: $border-radius;

.space-output {
	border: medium solid darkblue;
	background-color: black;
	border-radius: $border-radius;
	padding: 20px;
	cursor: pointer; // clickable spaces

	>.space-info-bar {
		background-color: rgb(200, 200, 255);
		border-radius: $border-radius;
		padding: 10px;
		cursor: default;
	}

	>.space-title-bar {
		background-color: rgb(100, 100, 200);
		border-radius: $border-radius;
		padding: 10px;
		cursor: default;

		.space-title-text {
			background-color: white;
			padding: 5px;
		}
	}

}
</style>
