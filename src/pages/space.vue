<template>
<div class="space-page page-width-xl">

	<return-to-top/>

	<space-loader :space-id="spaceId" include-parent-path include-tags>

		<template #default="{userAllowPin, userAllowPinSubspacesOnCreate}">

			<div @click.stop class="subspaces flex-column-lg">

				<div class="flex-row-md">
					<create-dropdown
						:parent-id="spaceId"
						:disabled="$store.getters.createDisabled"
						/>
					<subspaces-filter
						v-if="subspaces"
						@update:filter="f => filter = f"
						/>
				</div>

				<space-output
					v-for="s in subspaces"
					:space="s"
					:user-allow-pin="userAllowPin"
					:user-allow-pin-subspaces-on-create="userAllowPinSubspacesOnCreate"
					@toggle-pinned="$emit('toggle-pinned', s)"
				/>

			</div>

		</template>

	</space-loader>

</div>
</template>

<script>
import SpaceLoader from '@/widgets/space-loader.vue';
import SpaceOutput from '@/widgets/space-output.vue';
import CreateDropdown from '@/widgets/create-dropdown.vue';
import SubspacesFilter from '@/widgets/subspaces-filter.vue';

import {
	ajaxGet,
} from '@/utils/ajax.js';

import {
	SPACE_TYPES,
} from '@/const.js';

export default {
	components: {
		SpaceLoader,
		SpaceOutput,
		CreateDropdown,
		SubspacesFilter,
	},
	data() {
		return {
			loading: true,
			space: null,
			subspaces: [],
			filter: null,
			loadingSubspaces: false,
			skipLoadSubspaces: true,
		};
	},
	computed: {
		spaceId() {
			return this.$route.params.spaceId;
		},
		filterJSON() {
			return this.filter ? JSON.stringify(this.filter) : null;
		},
	},
	watch: {
		filter: {
			deep: true,
			handler() {
				this.loadSubspaces();
			},
		},
	},
	methods: {
		loadSubspaces(more = false) {
			this.loadingSubspaces = true;

			if (!more) {
				this.subspaces = [];
			}

			ajaxGet('/ajax/subspaces', {
				parentId: this.spaceId,
				offset: this.subspaces.length,
				limit: this.$store.getters.maxPageLimit,
				filter: this.filterJSON,
			}).then(response => {
				this.subspaces = this.subspaces.concat(response);
			});
		},
	},
};
</script>
