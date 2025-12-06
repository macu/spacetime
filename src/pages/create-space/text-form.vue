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

	<form-field title="Text" required>
		<div ref="textBodyWrapper">
			<el-input
				type="textarea"
				v-model="text"
				:maxlength="$store.getters.textMaxLength"
				show-word-limit
				:autosize="{minRows: 3}"
				:disabled="posting"
				/>
		</div>
	</form-field>

	<form-actions>
		<el-checkbox v-model="saveRecording" :disabled="posting">
			Include typing record for replay
		</el-checkbox>
	</form-actions>

	<form-actions v-if="saveRecording">
		<el-button v-if="previewing" @click="previewing = false">
			Close preview
		</el-button>
		<el-button v-else @click="preview()" :disabled="recording.length === 0">
			Preview
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
		createDisabled() {
			return this.posting || !this.text.trim();
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
				event: 'change',
				ts: Date.now() - this.getStartedAt(),
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
		getStartedAt() {
			if (!this.startedAt) {
				this.startedAt = Date.now();
			}
			return this.startedAt;
		},
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

			const textarea = event.target;
			const selectionStart = textarea.selectionStart;
			const selectionEnd = textarea.selectionEnd;
			const timestamp = Date.now() - this.getStartedAt();
			this.recording.push({
				event: selectionStart === selectionEnd ? 'cursor' : 'select',
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
				saveRecording: this.saveRecording,
				startedAt: this.startedAt,
				recording: this.saveRecording ? this.recording : [],
			});
		},
	},
};
</script>
