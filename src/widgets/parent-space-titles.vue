<template>
<div class="parent-space-titles flex-row">

	<loading-message v-if="loading"/>

	<div
		v-for="p in parents"
		:key="p.id"
		class="parent-space"
		@click="gotoSpace(p)">

		<div @click.stop
			class="space-titles flex-row-md nowrap horizontal-scroll">
			<template v-if="p.userTitles">
				<space-title
					v-for="t in p.userTitles"
					:key="t.id"
					:space="t"
					@check-in="checkinUserTitle(p, t)"
					/>
			</template>
			<template v-if="p.topTitles">
				<space-title
					v-for="t in p.topTitles"
					:key="t.id"
					:space="t"
					@check-in="checkinTopTitle(p, t)"
					/>
			</template>
		</div>

	</div>

</div>
</template>

<script>
import SpaceTitle from '@/widgets/space-title.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

export default {
	components: {
		SpaceTitle,
	},
	props: {
		parentId: {
			type: Number,
			required: true,
		},
	},
	data() {
		return {
			loading: true,
			parents: [],
		};
	},
	mounted() {
		ajaxGet('/ajax/space/parent-titles', {
			parentId: this.parentId,
		}).then(response => {
			this.parents = response;
		}).finally(() => {
			this.loading = false;
		});
	},
	methods: {
		gotoSpace(space) {
			this.$router.push({
				name: 'space',
				params: {
					spaceId: space.id,
				},
			});
		},
		checkinUserTitle(parent, title) {
		},
		checkinTopTitle(parent, title) {
		},
	},
};
</script>
