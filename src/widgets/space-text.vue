<template>
<div class="space-text" @click.stop>

	<el-button v-if="showLoadRecording" type="primary" @click="loadRecording()">
		<material-icon icon="play_circle" />
		<span>Play recording</span>
	</el-button>
	<loading-message v-else-if="loading"/>

	<text-replay
		v-if="replayData"
		:recording="replayData"
		:final-text="textOutput"
		initial-playing
		inline
	/>
	<div v-else-if="!loading" v-text="textOutput"/>

</div>
</template>

<script>
import TextReplay from '@/widgets/text-replay.vue';
import {ajaxGet} from '@/utils/ajax.js';

export default {
	components: {
		TextReplay,
	},
	props: {
		space: {
			type: Object,
			required: true,
		},
	},
	data() {
		return {
			loading: false,
			replayData: null,
		};
	},
	computed: {
		textOutput() {
			return (this.space.text || '').trim();
		},
		hasRecording() {
			return this.space.hasRecording || false;
		},
		showLoadRecording() {
			return this.hasRecording && !this.loading && !this.replayData;
		},
	},
	methods: {
		loadRecording() {
			ajaxGet('/ajax/text-replay', {
				spaceId: this.space.id,
			}).then(recording => {
				this.replayData = recording || null;
			}).finally(() => {
				this.loading = false;
			});
		},
	},
};
</script>

<style lang="scss">
@import '@/styles/vars.scss';

.space-text {
	background-color: $text-bg-color;
	border: $text-border;
	color: $text-color;
	border-radius: $border-radius;
	box-shadow: $text-inner-drop-shadow;

	padding: 20px;

	font-size: 1.2em;
	white-space: pre-wrap;

	display: flex;
	flex-direction: column;
	gap: 1em;

	>.el-button {
		align-self: flex-start;
	}
}
</style>
