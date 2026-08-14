<template>
<div class="space-output">

	<div v-if="showPath && hasParentPath" @click.stop class="parent-path">
		<div
			v-for="p in space.parentPath"
			:key="p.id"
			@click.stop="gotoSpace(p, true)"
			class="flex-row-md" :class="{'clickable': gotoPathOnClick}">

			<material-icon icon="arrow_right_alt"/>

			<el-tooltip v-if="p.isPinned" content="Pinned by author" placement="top">
				<material-icon icon="keep"/>
			</el-tooltip>

			<space-type :space="p"/>

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

		<div class="space-header flex-row-md nowrap"
			:class="{'clickable': gotoSpaceOnClick}"
			@click.stop="gotoSpace()">
			<div class="flex-row-md">
				<el-tooltip v-if="space.isPinned" content="Pinned by author" placement="top">
					<material-icon icon="keep"/>
				</el-tooltip>
				<space-type :space="space" @click="gotoSpace()"/>
				<div v-if="space.label" class="space-label">
					<strong v-text="space.label"/>
				</div>
				<space-creator :space="space"/>
			</div>
			<template v-if="showReorder">
				<div class="flex-1"/>
				<material-icon class="drag-handle" icon="drag_indicator"/>
			</template>
		</div>

		<slot name="actions-area" :space="space" :context="context"/>

		<space-tag
			v-if="space.spaceType === SPACE_TYPES.TAG"
			:space="space"
			:show-actions="false"
			flat
			/>

		<space-text
			v-else-if="space.spaceType === SPACE_TYPES.TEXT"
			:space="space"
			/>

		<template v-else-if="space.spaceType === SPACE_TYPES.LINK">
			<space-output
				v-if="space.linkSpace"
				:space="space.linkSpace"
				show-path
				:goto-space-on-click="gotoSpaceOnClick"
				:goto-path-on-click="gotoPathOnClick"
			/>
			<el-alert v-else type="error" :closable="false">
				<p>This link points to a space that no longer exists.</p>
			</el-alert>
		</template>

		<div v-if="$slots.default" class="portal" @click.stop>
			<slot :context="context"/>
		</div>

	</div>

</div>
</template>

<script>
import SpaceType from './space-type.vue';
import SpaceCreator from './space-creator.vue';
import SpaceTag from './space-tag.vue';
import SpaceText from './space-text.vue';

import {
	ajaxPost,
} from '@/utils/ajax.js';

import {
	SPACE_TYPES,
} from '@/const.js';

import {
	alertSuccess,
} from '@/utils/notify.js';

export default {
	name: 'space-output', // recursive
	emits: [
		'set-pinned',
	],
	components: {
		SpaceType,
		SpaceCreator,
		SpaceTag,
		SpaceText,
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
		gotoSpaceOnClick: {
			type: Boolean,
			default: true,
		},
		gotoPathOnClick: {
			type: Boolean,
			default: true,
		},
		context: {
			type: Object,
			required: false,
		},
		subSpace: {
			type: Boolean,
			default: false,
		},
		showReorder: {
			type: Boolean,
			default: false,
		},
	},
	computed: {
		SPACE_TYPES() {
			return SPACE_TYPES;
		},
		hasParentPath() {
			return !!this.space.parentPath && this.space.parentPath.length > 0;
		},
	},
	methods: {
		gotoSpace(s = null, path = false) {
			if (!this.gotoSpaceOnClick) {
				return;
			}
			if (!path || this.gotoPathOnClick) {
				this.$router.push({
					name: 'space',
					params: {
						spaceId: s ? s.id : this.space.id,
					},
				});
			}
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
			&.clickable {
				cursor: pointer;
			}
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
		padding: 20px 20px;

		>.space-header {
			&.clickable {
				cursor: pointer;
			}

			.drag-handle {
				cursor: ns-resize;
			}

			.space-label {
				padding-left: 10px;
				font-weight: bold;
				font-size: 120%;
			}
		}

		>.space-tag {
			padding: 40px;
			cursor: default;
		}
		>.space-text {
			padding: 40px;
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
		>.space-tag, >.space-text {
			padding: 20px;
		}
		>.portal {
			padding: 20px 10px;
		}
	}
}
</style>
