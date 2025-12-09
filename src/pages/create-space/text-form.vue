<template>
<form-layout title="Create text">

	<form-field title="Title">
		<el-input
			v-model="title"
			:maxlength="$store.getters.titleMaxLength"
			show-word-limit
			size="large"
			:disabled="posting"
			/>
	</form-field>

	<form-actions>
		<el-checkbox v-model="saveRecording" :disabled="disableRecordingOption">
			Include typing record for replay
		</el-checkbox>
	</form-actions>

	<form-field title="Text" required>
		<el-progress
			v-if="saveRecording"
			:percentage="(recording.length / $store.getters.nakedTextMaxDeltas) * 100"
			:format="pct => Math.floor(pct) + '%'"
			:stroke-width="24"
			show-text
			/>
		<div ref="textBodyWrapper">
			<el-input
				type="textarea"
				v-model="text"
				:maxlength="$store.getters.textMaxLength"
				show-word-limit
				:autosize="{minRows: 3}"
				:disabled="disableTextarea"
				/>
		</div>
	</form-field>

	<form-actions v-if="saveRecording">
		<el-button v-if="previewing" @click="previewing = false">
			Close preview
		</el-button>
		<el-button v-else @click="preview()" :disabled="recording.length === 0">
			Preview recording
		</el-button>
	</form-actions>

	<text-replay
		v-if="previewing"
		:recording="recording"
		:final-text="text"
		:initial-playing="true"
		/>

	<form-actions>
		<el-button @click="submit()" type="primary" :disabled="createDisabled">
			Create
		</el-button>
	</form-actions>

</form-layout>
</template>

<script>
import TextReplay from '@/widgets/text-replay.vue';

export default {
	components: {
		TextReplay,
	},
	props: {
		posting: {
			type: Boolean,
			default: false,
		},
		initialSaveRecording: {
			type: Boolean,
			default: false,
		},
	},
	data() {
		return {
			title: '',
			text: '',
			saveRecording: this.initialSaveRecording,
			startedAt: null,
			recording: [],

			previewing: false,
		};
	},
	computed: {
		recordingMaxed() {
			return this.recording.length > this.$store.getters.nakedTextMaxDeltas;
		},
		createDisabled() {
			return this.posting || !this.text.trim();
		},
		disableRecordingOption() {
			return this.posting || (!this.saveRecording && this.recordingMaxed);
		},
		disableTextarea() {
			return this.posting || (this.saveRecording && this.recordingMaxed);
		},
	},
	watch: {
		text(newValue, oldValue) {
			this.previewing = false;

			if (newValue === oldValue) {
				return;
			}
			if (newValue.length > this.$store.getters.textMaxLength) {
				// Should be safe using maxlength on input
				newValue = newValue.slice(0, this.$store.getters.textMaxLength);
			}

			// Ensure the first timestamp is 0
			let timestamp;
			if (this.startedAt) {
				timestamp = Date.now() - this.startedAt;
			} else {
				this.startedAt = Date.now();
				timestamp = 0;
			}

			// TODO Use selection for accurate positioning in repeated strings

			// Find start of delta in old value
			let changeStart = 0;
			while (
				changeStart < newValue.length &&
				changeStart < oldValue.length &&
				newValue[changeStart] === oldValue[changeStart]
			) {
				changeStart++;
			}

			// Find end of delta in old value
			// Search from end of both strings
			let oldValueChangeEnd = oldValue.length - 1;
			let newValueEndIndex = newValue.length - 1;
			while (
				oldValueChangeEnd >= changeStart &&
				oldValueChangeEnd >= 0 &&
				newValueEndIndex >= 0 &&
				newValue[newValueEndIndex] === oldValue[oldValueChangeEnd]
			) {
				oldValueChangeEnd--;
				newValueEndIndex--;
			}

			// Add delta
			this.recording.push({
				et: 'change',
				ts: timestamp,
				ss: changeStart, // selection applied before delta
				se: oldValueChangeEnd + 1,
				t: newValue.slice(changeStart, newValueEndIndex + 1),
			});
		},
	},
	mounted() {
		let textarea = this.$refs.textBodyWrapper.querySelector('textarea');
		if (textarea) {
			textarea.focus();
			// Add selection change
			textarea.addEventListener('select', this.onSelect);
			// Add cursor position change
			textarea.addEventListener('keyup', this.onSelect);
			textarea.addEventListener('click', this.onSelect);
			textarea.addEventListener('focus', this.onSelect);
		}
	},
	beforeUnmount() {
		let textarea = this.$refs.textBodyWrapper.querySelector('textarea');
		if (textarea) {
			textarea.removeEventListener('select', this.onSelect);
			textarea.removeEventListener('keyup', this.onSelect);
			textarea.removeEventListener('click', this.onSelect);
			textarea.removeEventListener('focus', this.onSelect);
		}
	},
	methods: {
		onSelect(event) {
			// Track cursor navigation
			if (event.type === 'keyup') {
				const allowedKeys = [
					'ArrowLeft',
					'ArrowRight',
					'ArrowUp',
					'ArrowDown',
					'Home',
					'End',
					'PageUp',
					'PageDown',
					// Try to catch combo movements
					'Control',
					'Alt',
					'Meta',
					'Shift',
				];
				if (!allowedKeys.includes(event.key)) {
					return;
				}
			}

			// Ensure the first timestamp is 0
			let timestamp;
			if (this.startedAt) {
				timestamp = Date.now() - this.startedAt;
			} else {
				this.startedAt = Date.now();
				timestamp = 0;
			}

			const textarea = event.target;
			const selectionStart = textarea.selectionStart;
			const selectionEnd = textarea.selectionEnd;

			if (selectionStart === selectionEnd &&
				this.recording.length > 0 &&
				this.recording[this.recording.length - 1].et === 'cursor' &&
				this.recording[this.recording.length - 1].ss === selectionStart
			) {
				// No change in cursor position
				return;
			}

			this.recording.push({
				et: selectionStart === selectionEnd ? 'cursor' : 'select',
				ts: timestamp,
				ss: selectionStart,
				se: selectionEnd,
			});
		},
		preview() {
			this.previewing = true;
		},
		submit() {
			if (this.createDisabled) {
				return;
			}
			this.$emit('submit', {
				title: this.title.trim(),
				text: this.text.trim(),
				recording: this.saveRecording ? JSON.stringify(this.recording) : null,
				startedAt: this.startedAt,
			});
		},
	},
};
</script>
