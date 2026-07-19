<template>
<div class="form-actions flex-row-md">
	<slot/>
</div>
</template>

<style lang="scss">
@import '@/styles/vars.scss';

.form-actions {
	padding: 10px;
	border: thin solid black;
	border-radius: $border-radius;
}
</style>
