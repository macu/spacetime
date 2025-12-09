<template>
<div class="naked-text-replay" :class="{ 'inline': inline }">
	<div v-if="currentSpans">
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
		:text="progressDisplay"
		/>
	<div class="actions flex-row">
		<el-button v-if="playing" @click="stop()" type="warning" size="small">
			Stop
		</el-button>
		<el-button v-else @click="playing = true" size="small" type="primary">
			Play
		</el-button>
		<el-button v-if="showSkipAhead" @click="skipToNext()" type="primary" size="small">
			Skip to next event ({{displayTimeToNextEvent}})
		</el-button>
	</div>
</div>
</template>

<script>
import {formatTimestamp} from '@/utils/time.js';

const PADDING_START = 1000; // 1 second at start before first event
const PADDING_END = 5000; // 5 seconds at end before stop
const SKIP_AHEAD_THRESHOLD = 3000; // 10 seconds ahead to skip to end

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
		inline: {
			type: Boolean,
			default: false,
		},
	},
	setup() {
		return {
			// non-reactive properties
			timeout: null,
			currentTimeInterval: null,
		};
	},
	data() {
		return {
			playing: this.initialPlaying,
			currentIndex: 0,
			currentText: '',
			currentEvent: null,
			startPlayingAt: null,
			currentTime: 0,
		};
	},
	computed: {
		currentSpans() {
			if (!this.playing) {
				return null;
			}

			if (!this.currentEvent) {
				// Padding time before first event
				return [
					{ text: '', classes: ['cursor'] },
				];
			}

			let before = this.currentText.slice(0, this.currentEvent.ss);
			let inserted = this.currentEvent.t || '';
			let after = this.currentText.slice(this.currentEvent.se + inserted.length);

			if (this.currentEvent.et === 'change') {
				// This is a text insertion/deletion
				// Display cursor after newly inserted text

				return [
					{ text: before, classes: [] },
					{ text: inserted, classes: ['inserted-text'] },
					{ text: '', classes: ['cursor'] },
					{ text: after, classes: [] },
				];

			} else if (this.currentEvent.et === 'cursor') {
				// Cursor placement

				return [
					{ text: before, classes: [] },
					{ text: '', classes: ['cursor'] },
					{ text: after, classes: [] },
				];

			} else if (this.currentEvent.et === 'select') {
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
		progressDisplay() {
			return `${this.currentIndex} / ${this.recording.length}`;
		},
		timeToNextEvent() {
			if (!this.playing || !this.currentEvent ||
				this.currentIndex >= (this.recording.length - 1)
			) {
				return null;
			}
			return this.startPlayingAt +
				this.recording[this.currentIndex + 1].ts - this.currentTime;
		},
		showSkipAhead() {
			if (!this.timeToNextEvent) {
				return false;
			}
			return this.timeToNextEvent > SKIP_AHEAD_THRESHOLD;
		},
		displayTimeToNextEvent() {
			return this.timeToNextEvent
				? `in ${formatTimestamp(this.timeToNextEvent)}`
				: null;
		},
	},
	watch: {
		playing: {
			immediate: true,
			handler(newVal) {
				if (newVal) {
					this.playNext();
				} else {
					if (this.timeout) {
						clearTimeout(this.timeout);
						this.timeout = null;
					}
					if (this.currentTimeInterval) {
						clearInterval(this.currentTimeInterval);
						this.currentTimeInterval = null;
					}
				}
			},
		},
	},
	beforeUnmount() {
		if (this.timeout) {
			clearTimeout(this.timeout);
			this.timeout = null;
		}
		if (this.currentTimeInterval) {
			clearInterval(this.currentTimeInterval);
			this.currentTimeInterval = null;
		}
	},
	methods: {
		skipToNext() {
			if (!this.playing || this.currentIndex >= (this.recording.length - 1)) {
				return;
			}

			this.startPlayingAt = (Date.now() - this.recording[this.currentIndex + 1].ts);

			this.playNext();
		},
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
				this.startPlayingAt = Date.now() + PADDING_START;
				this.timeout = setTimeout(() => {
					this.playNext();
				}, PADDING_START);
				this.currentTimeInterval = setInterval(() => {
					this.currentTime = Date.now();
				}, 1000);
				return;
			}

			this.currentEvent = this.recording[this.currentIndex];

			if (this.currentEvent.et === 'change') {
				// Apply diff to currentText
				let before = this.currentText.slice(0, this.currentEvent.ss);
				let after = this.currentText.slice(this.currentEvent.se);
				this.currentText = before + (this.currentEvent.t || '') + after;
			}

			this.currentIndex++;

			if (this.currentIndex < this.recording.length) {
				let nextEvent = this.recording[this.currentIndex];
				let delay = this.startPlayingAt + nextEvent.ts - Date.now();
				this.timeout = setTimeout(() => {
					this.playNext();
				}, Math.max(0, delay));
			} else {
				// Finished
				this.timeout = setTimeout(() => {
					this.stop();
				}, PADDING_END);
			}

		},
		stop() {
			this.playing = false;
			this.currentEvent = null;
			this.startPlayingAt = null;
			this.currentText = '';
			this.currentIndex = 0;
			if (this.timeout) {
				clearTimeout(this.timeout);
				this.timeout = null;
			}
			if (this.currentTimeInterval) {
				clearInterval(this.currentTimeInterval);
				this.currentTimeInterval = null;
			}
			this.currentTime = 0;
			this.$emit('finished');
		},
	},
};
</script>

<style lang="scss">
.naked-text-replay {

	&:not(.inline) {
		padding: 1em;
		background: white;
		border: 1px solid #ccc;
	}

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
				animation: highlight-insert 2s ease-out forwards;
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
