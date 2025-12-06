<template>
<div class="naked-text-replay">
	<div v-if="playing && currentSpans">
		<span
			v-for="s in currentSpans"
			v-text="s.text"
			:class="s.classes"
			/>
	</div>
	<div v-else v-text="finalText"/>
	<el-progress
		v-if="playing && recording.length > 0"
		:percentage="percentage"
		:show-text="false"
		/>
	<div class="actions flex-row">
		<el-button @click="playing = !playing" size="small">
			{{ playing ? 'Pause' : 'Play' }}
		</el-button>
		<el-button @click="restart()" size="small">
			Restart
		</el-button>
	</div>
</div>
</template>

<script>
export default {
	emits: [
		'finished',
	],
	props: {
		recording: {
			type: Array,
			required: true,
		},
		finalText: {
			type: String,
			required: true,
		},
		initialPlaying: {
			type: Boolean,
			default: false,
		},
	},
	data() {
		return {
			playing: this.initialPlaying,
			currentIndex: 0,
			currentText: '',
			currentEvent: null,
			startPlayingAt: null,
		};
	},
	computed: {
		currentSpans() {
			if (!this.playing || !this.currentEvent) {
				return null;
			}

			let before = this.currentText.slice(0, this.currentEvent.ss);
			let inserted = this.currentEvent.t || '';
			let after = this.currentText.slice(this.currentEvent.se + inserted.length);

			if (this.currentEvent.event === 'change') {
				// This is a text insertion/deletion
				// Display cursor after newly inserted text

				return [
					{ text: before, classes: [] },
					{ text: inserted, classes: ['inserted-text'] },
					{ text: '', classes: ['cursor'] },
					{ text: after, classes: [] },
				];

			} else if (this.currentEvent.event === 'cursor') {
				// Cursor placement

				return [
					{ text: before, classes: [] },
					{ text: '', classes: ['cursor'] },
					{ text: after, classes: [] },
				];

			} else if (this.currentEvent.event === 'select') {
				// Text selection

				let selectedText = this.currentText.slice(
					this.currentEvent.ss,
					this.currentEvent.se,
				);

				return [
					{ text: before, classes: [] },
					{ text: selectedText, classes: ['selected-text'] },
					{ text: after, classes: [] },
				];

			}

			return [
				{ text: this.currentText, classes: [] },
			];
		},
		percentage() {
			return (this.currentIndex / this.recording.length) * 100;
		},
	},
	watch: {
		playing: {
			immediate: true,
			handler(newVal) {
				if (newVal) {
					this.playNext();
				}
			},
		},
	},
	methods: {
		playNext() {
			if (!this.playing) {
				return;
			}
			if (this.currentIndex >= this.recording.length) {
				this.stop();
				this.$emit('finished');
				return;
			}

			if (this.startPlayingAt === null) {
				this.startPlayingAt = Date.now();
			}

			this.currentEvent = this.recording[this.currentIndex];

			if (this.currentEvent.event === 'change') {
				// Apply diff to currentText
				let before = this.currentText.slice(0, this.currentEvent.ss);
				let after = this.currentText.slice(this.currentEvent.se);
				this.currentText = before + (this.currentEvent.t || '') + after;
			}

			this.currentIndex++;

			if (this.currentIndex < this.recording.length) {
				let nextEvent = this.recording[this.currentIndex];
				let delay = this.startPlayingAt + nextEvent.ts - Date.now();
				setTimeout(() => {
					this.playNext();
				}, Math.max(0, delay));
			} else {
				// Finished
				this.playing = false;
			}

		},
		stop() {
			this.playing = false;
			this.currentEvent = null;
			this.startPlayingAt = null;
			this.currentText = '';
			this.currentIndex = 0;
		},
		restart() {
			this.currentIndex = 0;
			this.currentText = '';
			this.currentEvent = null;
			this.startPlayingAt = null;
			this.playing = true;
		},
	},
};
</script>

<style lang="scss">
.naked-text-replay {

	padding: 1em;
	background: white;
	border: 1px solid #ccc;

	display: flex;
	flex-direction: column;
	gap: 1em;

	>div:not([class]) {
		white-space: pre-wrap;
		word-break: break-word;
		>span {
			&.inserted-text {
				background-color: #d4f8d4;
				// animate background to white and stop
				animation: highlight-insert 2s ease-out;
			}
			&.selected-text {
				background-color: #add8ff;
			}
			&.cursor {
				display: inline-block;
				width: 1px;
				background-color: black;
				height: 1em;
				animation: blink 1s steps(1) infinite;
				vertical-align: bottom;
			}
		}
	}

}

@keyframes blink {
	0% { opacity: 1; }
	50% { opacity: 0; }
	100% { opacity: 1; }
}

@keyframes highlight-insert {
	0% { background-color: #d4f8d4; }
	100% { background-color: white; }
}
</style>
