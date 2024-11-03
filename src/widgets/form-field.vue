<template>
<div class="form-layout-form-field flex-column">
	<div v-if="showTitle" class="flex-row">
		<h3 class="flex-1"><slot name="title">{{title}}</slot></h3>
		<em v-if="required">Required</em>
	</div>
	<div class="form-field-body flex-column">
		<slot><span v-if="showValue" class="form-value" v-text="value"/></slot>
	</div>
</div>
</template>

<script>
export default {
	props: {
		title: {
			type: String,
			required: false,
		},
		value: {
			type: [String, Number],
			required: false,
		},
		required: {
			type: Boolean,
			default: false,
		},
	},
	computed: {
		showTitle() {
			return !!this.title || this.$slots.title !== undefined;
		},
		showValue() {
			return this.value !== null;
		},
	},
};
</script>

<style lang="scss">
.form-layout-form-field {
	>h3 {
		font-weight: bold;
	}
	>.form-field-body {
		padding: 20px;
		background-color: #eee;
		border-radius: 8px;
		>.form-value {
			font-size: larger;
		}
	}
}
</style>
