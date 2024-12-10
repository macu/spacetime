<template>
<div class="space-output flex-column">

	<div class="space-info-area flex-row-md">
		<checkin-button :space="space"/>
		<space-type :type="space.spaceType"/>
	</div>

	<div v-if="hasTitle" class="space-title-area flex-row-md">
		<space-title v-if="!!space.lastUserTitle"
			:title-space="space.lastUserTitle"
			/>
		<space-title v-for="title in space.topTitles"
			:title-space="title"
			/>
	</div>

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
		/>

</div>
</template>

<script>
import CheckinButton from './checkin-button.vue';
import SpaceType from './space-type.vue';
import SpaceTitle from './space-title.vue';
import SpaceText from './space-text.vue';

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
	},
	props: {
		space: {
			type: Object,
			required: true,
		},
	},
	computed: {
		SPACE_TYPES() {
			return SPACE_TYPES;
		},
		hasTitle() {
			return !!this.space.lastUserTitle ||
				(!!this.space.topTitles && this.space.topTitles.length > 0);
		},
	},
	methods: {

	},
};
</script>

<style lang="scss">
.space-output {
	border: thin solid darkblue;
	border-radius: 12px;
	padding: 20px;


}
</style>
